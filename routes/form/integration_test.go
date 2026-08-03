package form_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/middleware"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	m "github.com/imcrazytwkr/formdrain/models/http"
	fcr "github.com/imcrazytwkr/formdrain/repositories/form_config"
	frr "github.com/imcrazytwkr/formdrain/repositories/form_response"
	scr "github.com/imcrazytwkr/formdrain/repositories/site_config"
	"github.com/imcrazytwkr/formdrain/routes/form"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

type fakeCaptcha struct {
	err   error
	calls int
}

func (f *fakeCaptcha) Validate(
	_ context.Context,
	_ fc.CaptchaType,
	_ string,
	_ string,
	_ netip.Addr,
) error {
	f.calls++
	return f.err
}

type recordingNotifier struct {
	mu    sync.Mutex
	calls int
	last  map[string]any
	err   error
	done  chan struct{}
}

func newRecordingNotifier() *recordingNotifier {
	return &recordingNotifier{done: make(chan struct{}, 8)}
}

func (n *recordingNotifier) Send(_ context.Context, _ fc.NotifiersConfig, formData map[string]any) error {
	n.mu.Lock()
	n.calls++
	n.last = formData
	err := n.err
	n.mu.Unlock()

	select {
	case n.done <- struct{}{}:
	default:
	}
	return err
}

func (n *recordingNotifier) SendAsync(ctx context.Context, config fc.NotifiersConfig, formData map[string]any) {
	go func() {
		_ = n.Send(ctx, config, formData)
	}()
}

func (n *recordingNotifier) waitCall(t *testing.T) {
	t.Helper()
	select {
	case <-n.done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for notifier")
	}
}

func (n *recordingNotifier) snapshot() (int, map[string]any) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.calls, n.last
}

type harness struct {
	db       *sql.DB
	handler  http.Handler
	captcha  *fakeCaptcha
	notifier *recordingNotifier
}

func newHarness(t *testing.T) *harness {
	t.Helper()

	db := testutil.OpenSqlite(t)
	captcha := &fakeCaptcha{}
	notifier := newRecordingNotifier()

	router := chi.NewRouter()
	router.Use(middleware.ResponseFormatParser(m.ContentTypeHTML, m.ContentTypeJSON))
	router.Route("/form", func(r chi.Router) {
		form.NewFormRouter(
			fcr.NewSqliteFormConfigRepository(db),
			scr.NewSqliteSiteConfigRepository(db),
			frr.NewSqliteFormResponseRepository(db),
			captcha,
			notifier,
		).Router(r)
	})

	return &harness{
		db:       db,
		handler:  router,
		captcha:  captcha,
		notifier: notifier,
	}
}

func (h *harness) seed(t *testing.T, formID int64, redirectTo string) {
	t.Helper()

	_, err := h.db.Exec(`INSERT OR IGNORE INTO sites (id, hostname) VALUES (1, 'example.com')`)
	if err != nil {
		t.Fatalf("seed site: %v", err)
	}

	_, err = h.db.Exec(`DELETE FROM forms WHERE id = ?`, formID)
	if err != nil {
		t.Fatalf("seed delete form: %v", err)
	}

	var redirect any
	if redirectTo != "" {
		redirect = redirectTo
	}

	_, err = h.db.Exec(`
		INSERT INTO forms (id, site_id, captcha_type, redirect_to, field_schema, schema_version, notifiers)
		VALUES (
			?,
			1,
			'hcaptcha',
			?,
			'{"version":1,"fields":[{"name":"email","type":"string","required":true}]}',
			1,
			'{"discord":null,"brevo":null}'
		)
	`, formID, redirect)
	if err != nil {
		t.Fatalf("seed form: %v", err)
	}
}

func (h *harness) postJSON(t *testing.T, formID string, body string, hdr http.Header) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/form/"+formID+"/", strings.NewReader(body))
	req.RemoteAddr = "203.0.113.10:1234"
	req.Host = "api.example.com"
	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJson)
	req.Header.Set(constants.HeaderOrigin, "https://example.com")
	for k, vals := range hdr {
		for _, v := range vals {
			req.Header.Set(k, v)
		}
	}
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	return w
}

func (h *harness) postForm(t *testing.T, formID string, values url.Values, hdr http.Header) *httptest.ResponseRecorder {
	t.Helper()
	encoded := values.Encode()
	req := httptest.NewRequest(http.MethodPost, "/form/"+formID+"/", strings.NewReader(encoded))
	req.RemoteAddr = "203.0.113.10:1234"
	req.Host = "api.example.com"
	req.Header.Set(constants.HeaderContentType, constants.ContentTypeForm)
	req.Header.Set(constants.HeaderOrigin, "https://example.com")
	for k, vals := range hdr {
		for _, v := range vals {
			req.Header.Set(k, v)
		}
	}
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	return w
}

func (h *harness) responseCount(t *testing.T, formID int64) int {
	t.Helper()
	var n int
	err := h.db.QueryRow(`SELECT COUNT(*) FROM form_responses WHERE form_id = ?`, formID).Scan(&n)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	return n
}

func TestCreate_JSONHappyPath(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	hdr := http.Header{}
	hdr.Set(constants.HeaderAccept, m.ContentTypeJSON.String())
	w := h.postJSON(t, "10", `{"email":"a@b.c","h-captcha":"tok"}`, hdr)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body = %q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Header().Get(constants.HeaderContentType), "application/json") {
		t.Fatalf("content-type = %q", w.Header().Get(constants.HeaderContentType))
	}
	if h.captcha.calls != 1 {
		t.Fatalf("captcha calls = %d", h.captcha.calls)
	}
	h.notifier.waitCall(t)
	calls, last := h.notifier.snapshot()
	if calls != 1 {
		t.Fatalf("notifier calls = %d", calls)
	}
	if _, ok := last["h-captcha"]; ok {
		t.Fatalf("captcha token leaked into payload: %#v", last)
	}
	if last["email"] != "a@b.c" {
		t.Fatalf("payload = %#v", last)
	}
	if h.responseCount(t, 10) != 1 {
		t.Fatal("expected one stored response")
	}

	var clientIP string
	var payload string
	err := h.db.QueryRow(`SELECT client_ip, payload FROM form_responses WHERE form_id = 10`).Scan(&clientIP, &payload)
	if err != nil {
		t.Fatalf("select: %v", err)
	}
	if clientIP != "203.0.113.10" {
		t.Fatalf("client_ip = %q", clientIP)
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(payload), &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["email"] != "a@b.c" {
		t.Fatalf("stored payload = %#v", parsed)
	}
}

func TestCreate_FormRedirect(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 11, "https://example.com/thanks")

	values := url.Values{}
	values.Set("email", "a@b.c")
	values.Set("h-captcha", "tok")
	w := h.postForm(t, "11", values, nil)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body = %q", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); loc != "https://example.com/thanks" {
		t.Fatalf("Location = %q", loc)
	}
	if h.responseCount(t, 11) != 1 {
		t.Fatal("expected stored response")
	}
}

func TestCreate_FormNotFound(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	cases := []string{"999", "abc", "0", "-1"}
	for _, formID := range cases {
		t.Run(formID, func(t *testing.T) {
			w := h.postJSON(t, formID, `{"email":"a@b.c"}`, nil)
			if w.Code != http.StatusNotFound {
				t.Fatalf("status = %d", w.Code)
			}
		})
	}
}

func TestCreate_CORSOriginMismatch(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	hdr := http.Header{}
	hdr.Set(constants.HeaderOrigin, "https://evil.example.com")
	w := h.postJSON(t, "10", `{"email":"a@b.c"}`, hdr)
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d body = %q", w.Code, w.Body.String())
	}
	if h.responseCount(t, 10) != 0 {
		t.Fatal("should not store response")
	}
}

func TestCreate_CORSRefererMismatch(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	hdr := http.Header{}
	hdr.Set(constants.HeaderOrigin, "https://example.com")
	hdr.Set(constants.HeaderReferer, "https://evil.example.com/x")
	w := h.postJSON(t, "10", `{"email":"a@b.c"}`, hdr)
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d body = %q", w.Code, w.Body.String())
	}
}

func TestCreate_NoOriginSkipsHandler(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	req := httptest.NewRequest(http.MethodPost, "/form/10/", strings.NewReader(`{"email":"a@b.c"}`))
	req.RemoteAddr = "203.0.113.10:1234"
	req.Host = "api.example.com"
	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJson)
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if h.captcha.calls != 0 || h.responseCount(t, 10) != 0 {
		t.Fatalf("handler ran without Origin: captcha=%d rows=%d",
			h.captcha.calls, h.responseCount(t, 10))
	}
	if calls, _ := h.notifier.snapshot(); calls != 0 {
		t.Fatalf("handler ran without Origin: notify=%d", calls)
	}
}

func TestCreate_Preflight(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	req := httptest.NewRequest(http.MethodOptions, "/form/10/", nil)
	req.Host = "api.example.com"
	req.Header.Set(constants.HeaderOrigin, "https://example.com")
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "example.com" {
		t.Fatalf("ACAO = %q", got)
	}
	allow := w.Header().Values("Access-Control-Allow-Methods")
	joined := strings.Join(allow, ",")
	if !strings.Contains(joined, http.MethodPost) || !strings.Contains(joined, http.MethodOptions) {
		t.Fatalf("Allow-Methods = %v", allow)
	}
}

func TestCreate_CaptchaNotPassed(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")
	h.captcha.err = constants.ErrCaptchaNotPassed

	w := h.postJSON(t, "10", `{"email":"a@b.c","h-captcha":"tok"}`, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body = %q", w.Code, w.Body.String())
	}
	if h.responseCount(t, 10) != 0 {
		t.Fatal("should not store")
	}
}

func TestCreate_CaptchaUnexpectedError(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")
	h.captcha.err = errors.New("upstream down")

	w := h.postJSON(t, "10", `{"email":"a@b.c","h-captcha":"tok"}`, nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCreate_ValidationError(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	hdr := http.Header{}
	hdr.Set(constants.HeaderAccept, m.ContentTypeJSON.String())
	w := h.postJSON(t, "10", `{"email":"a@b.c","extra":true}`, hdr)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body = %q", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json: %v body=%q", err, w.Body.String())
	}
	if body["status"] != float64(http.StatusBadRequest) {
		t.Fatalf("body = %#v", body)
	}
	if h.responseCount(t, 10) != 0 {
		t.Fatal("should not store")
	}
}

func TestCreate_UnsupportedContentType(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	req := httptest.NewRequest(http.MethodPost, "/form/10/", bytes.NewReader([]byte("x")))
	req.RemoteAddr = "203.0.113.10:1234"
	req.Host = "api.example.com"
	req.Header.Set(constants.HeaderContentType, "text/plain")
	req.Header.Set(constants.HeaderOrigin, "https://example.com")
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCreate_BodyTooLarge(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	body := `{"email":"` + strings.Repeat("a", 5000) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/form/10/", strings.NewReader(body))
	req.ContentLength = int64(len(body))
	req.RemoteAddr = "203.0.113.10:1234"
	req.Host = "api.example.com"
	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJson)
	req.Header.Set(constants.HeaderOrigin, "https://example.com")
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d body = %q", w.Code, w.Body.String())
	}
}

func TestCreate_NotifierErrorStillSucceeds(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")
	h.notifier.err = errors.New("discord down")

	hdr := http.Header{}
	hdr.Set(constants.HeaderAccept, m.ContentTypeJSON.String())
	w := h.postJSON(t, "10", `{"email":"a@b.c","h-captcha":"tok"}`, hdr)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	h.notifier.waitCall(t)
	if calls, _ := h.notifier.snapshot(); calls != 1 {
		t.Fatalf("notifier calls = %d", calls)
	}
	if h.responseCount(t, 10) != 1 {
		t.Fatal("expected stored response")
	}
}

func TestCreate_JSONBodyInfersResponseFormat(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	// No Accept header: ContentTypeParser should set response CT from body CT.
	w := h.postJSON(t, "10", `{"email":"a@b.c","h-captcha":"tok"}`, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body = %q", w.Code, w.Body.String())
	}
	ct := w.Header().Get(constants.HeaderContentType)
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("content-type = %q", ct)
	}
	if !json.Valid(w.Body.Bytes()) {
		t.Fatalf("expected JSON body, got %q", w.Body.String())
	}
}

func TestCreate_CustomCaptchaField(t *testing.T) {
	h := newHarness(t)
	h.seed(t, 10, "")

	_, err := h.db.Exec(`UPDATE forms SET captcha_field = ? WHERE id = 10`, "cf-turnstile-response")
	if err != nil {
		t.Fatalf("set captcha_field: %v", err)
	}

	hdr := http.Header{}
	hdr.Set(constants.HeaderAccept, m.ContentTypeJSON.String())
	w := h.postJSON(t, "10", `{"email":"a@b.c","cf-turnstile-response":"tok"}`, hdr)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body = %q", w.Code, w.Body.String())
	}
	if h.captcha.calls != 1 {
		t.Fatalf("captcha calls = %d", h.captcha.calls)
	}
	h.notifier.waitCall(t)
	_, last := h.notifier.snapshot()
	if _, ok := last["cf-turnstile-response"]; ok {
		t.Fatalf("captcha token leaked into payload: %#v", last)
	}
	if last["email"] != "a@b.c" {
		t.Fatalf("payload = %#v", last)
	}
}
