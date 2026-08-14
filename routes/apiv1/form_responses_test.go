package apiv1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/models/form_response"
	frr "github.com/imcrazytwkr/formdrain/repositories/form_response"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestListFormResponses_OK(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(1, "owned.example")
	testutil.SeedForm(t, h.db, 10, 1)
	seedFormResponses(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?limit=2", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.FormResponseList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 2 {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.Items[0].Id.String() != "00000000-0000-4000-8000-00000000000a" {
		t.Fatalf("first = %s", page.Items[0].Id)
	}
	if page.Items[1].Id.String() != "00000000-0000-4000-8000-00000000000b" {
		t.Fatalf("second = %s", page.Items[1].Id)
	}
	if page.Items[0].FormId != 10 || page.Items[0].Payload["n"] != "a" {
		t.Fatalf("mapped = %#v", page.Items[0])
	}
	if page.Items[0].ClientIp == nil || *page.Items[0].ClientIp != "203.0.113.10" {
		t.Fatalf("client_ip = %#v", page.Items[0].ClientIp)
	}
	if page.NextCursor == nil || len(*page.NextCursor) < 1 {
		t.Fatal("expected next_cursor")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?limit=2&cursor="+*page.NextCursor, nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page2 api.FormResponseList
	if err := json.NewDecoder(w.Body).Decode(&page2); err != nil {
		t.Fatal(err)
	}
	if len(page2.Items) != 2 {
		t.Fatalf("page2 = %#v", page2.Items)
	}
	if page2.Items[0].Id.String() != "00000000-0000-4000-8000-00000000000c" {
		t.Fatalf("page2 first = %s", page2.Items[0].Id)
	}
	if page2.Items[1].Id.String() != "00000000-0000-4000-8000-00000000000d" {
		t.Fatalf("page2 second = %s", page2.Items[1].Id)
	}
	if page2.Items[1].ClientIp != nil {
		t.Fatalf("expected null client_ip, got %#v", page2.Items[1].ClientIp)
	}
	if page2.NextCursor != nil {
		t.Fatalf("unexpected next_cursor %q", *page2.NextCursor)
	}
}

func TestListFormResponses_Empty(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(1, "owned.example")
	testutil.SeedForm(t, h.db, 10, 1)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.FormResponseList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 0 {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.NextCursor != nil {
		t.Fatalf("unexpected next_cursor %q", *page.NextCursor)
	}
}

func TestListFormResponses_NotFound(t *testing.T) {
	h := newApiV1Harness(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/999/responses", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestListFormResponses_OtherOwner(t *testing.T) {
	h := newApiV1Harness(t)
	testutil.SeedSite(t, h.db, 8, "other.example")
	testutil.SeedForm(t, h.db, 80, 8)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/80/responses", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestListFormResponses_Unauthorized(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(1, "owned.example")
	testutil.SeedForm(t, h.db, 10, 1)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestListFormResponses_BadCursor(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(1, "owned.example")
	testutil.SeedForm(t, h.db, 10, 1)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?cursor=not-valid", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestListFormResponses_SortByField(t *testing.T) {
	h := newApiV1Harness(t)
	seedSortableForm(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?sort=email", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	page := decodeFormResponseList(t, w)
	if len(page.Items) != 3 {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.Items[0].Payload["email"] != "a" || page.Items[1].Payload["email"] != "b" {
		t.Fatalf("order = %#v", page.Items)
	}
	if page.Items[2].Payload["email"] != nil {
		t.Fatalf("nulls last = %#v", page.Items[2])
	}
	if page.NextCursor != nil {
		t.Fatalf("unexpected next_cursor %q", *page.NextCursor)
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?sort=email&order=desc", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	page = decodeFormResponseList(t, w)
	if page.Items[0].Payload["email"] != "b" || page.Items[1].Payload["email"] != "a" {
		t.Fatalf("desc = %#v", page.Items)
	}
	if page.Items[2].Payload["email"] != nil {
		t.Fatalf("nulls last desc = %#v", page.Items[2])
	}
}

func TestListFormResponses_FieldCursor(t *testing.T) {
	h := newApiV1Harness(t)
	seedSortableForm(t, h)
	repo := frr.NewSqliteFormResponseRepository(h.db)
	if _, err := repo.SaveFormResponse(t.Context(), &form_response.FormResponse{
		Id: "00000000-0000-4000-8000-00000000000d", FormId: 10, SchemaVersion: 1,
		Payload: map[string]any{},
	}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?sort=email&limit=1", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	page := decodeFormResponseList(t, w)
	if len(page.Items) != 1 || page.Items[0].Payload["email"] != "a" {
		t.Fatalf("page1 = %#v", page.Items)
	}
	if page.NextCursor == nil {
		t.Fatal("expected next_cursor")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?sort=email&limit=1&cursor="+*page.NextCursor, nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	page = decodeFormResponseList(t, w)
	if len(page.Items) != 1 || page.Items[0].Payload["email"] != "b" {
		t.Fatalf("page2 = %#v", page.Items)
	}
	if page.NextCursor == nil {
		t.Fatal("expected next_cursor")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?sort=email&limit=1&cursor="+*page.NextCursor, nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	page = decodeFormResponseList(t, w)
	if len(page.Items) != 1 || page.Items[0].Payload["email"] != nil {
		t.Fatalf("page3 = %#v", page.Items)
	}
	if page.NextCursor == nil {
		t.Fatal("expected next_cursor in null group")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?sort=email&limit=1&cursor="+*page.NextCursor, nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	page = decodeFormResponseList(t, w)
	if len(page.Items) != 1 || page.Items[0].Payload["email"] != nil {
		t.Fatalf("page4 = %#v", page.Items)
	}
	if page.NextCursor != nil {
		t.Fatalf("unexpected next_cursor %q", *page.NextCursor)
	}
}

func TestListFormResponses_MismatchedCursor(t *testing.T) {
	h := newApiV1Harness(t)
	seedSortableForm(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?sort=email&limit=1", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.FormResponseList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if page.NextCursor == nil {
		t.Fatal("expected next_cursor")
	}
	emailCursor := *page.NextCursor

	cases := []string{
		"/api/v1/forms/10/responses?sort=age&cursor=" + emailCursor,
		"/api/v1/forms/10/responses?sort=email&order=desc&cursor=" + emailCursor,
		"/api/v1/forms/10/responses?cursor=" + emailCursor,
	}
	for _, path := range cases {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil)
		req.Header.Set(constants.HeaderAccept, "application/json")
		req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
		w := httptest.NewRecorder()
		h.handler.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("%s: status = %d body=%s", path, w.Code, w.Body.String())
		}
	}
}

func TestListFormResponses_SortNumberAndBoolean(t *testing.T) {
	h := newApiV1Harness(t)
	seedSortableForm(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?sort=age", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	page := decodeFormResponseList(t, w)
	if page.Items[0].Payload["age"] != float64(2) || page.Items[1].Payload["age"] != float64(10) {
		t.Fatalf("age = %#v", page.Items)
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10/responses?sort=ok&order=desc", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	page = decodeFormResponseList(t, w)
	if page.Items[0].Payload["ok"] != true || page.Items[1].Payload["ok"] != false {
		t.Fatalf("ok = %#v", page.Items)
	}
}

func TestListFormResponses_BadSort(t *testing.T) {
	h := newApiV1Harness(t)
	seedSortableForm(t, h)

	cases := []string{
		"/api/v1/forms/10/responses?sort=nope",
		"/api/v1/forms/10/responses?sort=tags",
		"/api/v1/forms/10/responses?order=up",
		"/api/v1/forms/10/responses?sort=email&cursor=not-valid",
	}
	for _, path := range cases {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil)
		req.Header.Set(constants.HeaderAccept, "application/json")
		req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
		w := httptest.NewRecorder()
		h.handler.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("%s: status = %d body=%s", path, w.Code, w.Body.String())
		}
	}
}

func seedFormResponses(t *testing.T, h *apiV1Harness) {
	t.Helper()

	repo := frr.NewSqliteFormResponseRepository(h.db)
	ip := netip.MustParseAddr("203.0.113.10")
	mustSave := func(resp *form_response.FormResponse) {
		t.Helper()
		if _, err := repo.SaveFormResponse(t.Context(), resp); err != nil {
			t.Fatal(err)
		}
	}

	mustSave(&form_response.FormResponse{
		Id: "00000000-0000-4000-8000-00000000000a", FormId: 10,
		CreatedAt:     time.Date(2026, 1, 1, 0, 0, 3, 0, time.UTC),
		SchemaVersion: 1, ClientIP: ip, Payload: map[string]any{"n": "a"},
	})
	mustSave(&form_response.FormResponse{
		Id: "00000000-0000-4000-8000-00000000000b", FormId: 10,
		CreatedAt:     time.Date(2026, 1, 1, 0, 0, 2, 0, time.UTC),
		SchemaVersion: 1, ClientIP: ip, Payload: map[string]any{"n": "b"},
	})
	mustSave(&form_response.FormResponse{
		Id: "00000000-0000-4000-8000-00000000000c", FormId: 10,
		CreatedAt:     time.Date(2026, 1, 1, 0, 0, 1, 0, time.UTC),
		SchemaVersion: 1, ClientIP: ip, Payload: map[string]any{"n": "c"},
	})
	mustSave(&form_response.FormResponse{
		Id: "00000000-0000-4000-8000-00000000000d", FormId: 10,
		CreatedAt:     time.Date(2026, 1, 1, 0, 0, 1, 0, time.UTC),
		SchemaVersion: 1, Payload: map[string]any{"n": "d"},
	})
}

func seedSortableForm(t *testing.T, h *apiV1Harness) {
	t.Helper()
	h.seedSite(1, "owned.example")
	_, err := h.db.Exec(`
		INSERT INTO forms (id, site_id, captcha_type, field_schema, schema_version, notifiers)
		VALUES (10, 1, 'hcaptcha', ?, 1, '{}');
	`, `{"fields":[{"name":"email","type":"string","required":false},{"name":"age","type":"number","required":false},{"name":"ok","type":"boolean","required":false},{"name":"tags","type":"array","required":false,"items":{"type":"string"}}]}`)
	if err != nil {
		t.Fatalf("seed form: %v", err)
	}

	repo := frr.NewSqliteFormResponseRepository(h.db)
	mustSave := func(id string, payload map[string]any) {
		t.Helper()
		if _, err := repo.SaveFormResponse(t.Context(), &form_response.FormResponse{
			Id: id, FormId: 10, SchemaVersion: 1, Payload: payload,
		}); err != nil {
			t.Fatal(err)
		}
	}
	mustSave("00000000-0000-4000-8000-00000000000a", map[string]any{"email": "b", "age": 10, "ok": true})
	mustSave("00000000-0000-4000-8000-00000000000b", map[string]any{"email": "a", "age": 2, "ok": false})
	mustSave("00000000-0000-4000-8000-00000000000c", map[string]any{})
}

func decodeFormResponseList(t *testing.T, w *httptest.ResponseRecorder) api.FormResponseList {
	t.Helper()
	var page api.FormResponseList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	return page
}
