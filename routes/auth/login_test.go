package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/imcrazytwkr/formdrain/models/account"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/repositories"
	ar "github.com/imcrazytwkr/formdrain/repositories/account"
	sr "github.com/imcrazytwkr/formdrain/repositories/session"
	"github.com/imcrazytwkr/formdrain/routes/auth"
	as "github.com/imcrazytwkr/formdrain/services/account"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func newAuthHandler(t *testing.T) (http.Handler, *account.Account, repositories.SessionRepository) {
	t.Helper()

	db := testutil.OpenSqlite(t)
	accountRepository := ar.NewSqliteAccountRepository(db)
	sessionRepository := sr.NewMemorySessionRepository()
	accountService := as.NewService(accountRepository)

	hash, err := as.HashPassword("correct-password")
	if err != nil {
		t.Fatal(err)
	}

	account := &account.Account{Email: "user@example.com", PasswordHash: hash}
	err = accountRepository.Create(t.Context(), account)
	if err != nil {
		t.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.ResponseFormatParser(m.ContentTypeHTML, m.ContentTypeJSON))
	router.Route("/auth", auth.NewAuthRouter(sessionRepository, accountService).Router)

	return router, account, sessionRepository
}

func loginSession(t *testing.T, h http.Handler) string {
	t.Helper()

	payload, _ := json.Marshal(map[string]string{
		"email":    "user@example.com",
		"password": "correct-password",
	})
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/login", bytes.NewReader(payload))
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("login status = %d body=%s", w.Code, w.Body.String())
	}
	cookies := w.Result().Cookies()
	if len(cookies) < 1 || cookies[0].Name != constants.CookieSession || cookies[0].Value == "" {
		t.Fatalf("cookies = %#v", cookies)
	}
	return cookies[0].Value
}

func TestLoginGet_HTMLOrigin(t *testing.T) {
	h, _, _ := newAuthHandler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/auth/login?redirect=/admin", nil)
	req.Header.Set(constants.HeaderAccept, "text/html")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, `name="origin" value="/admin"`) {
		t.Fatalf("missing origin hidden field: %s", body)
	}
}

func TestLoginGet_InvalidRedirectBecomesRoot(t *testing.T) {
	h, _, _ := newAuthHandler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/auth/login?redirect=https://evil.example/", nil)
	req.Header.Set(constants.HeaderAccept, "text/html")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `name="origin" value="/"`) {
		t.Fatalf("body = %s", w.Body.String())
	}
}

func TestLoginGet_JSONNotAcceptable(t *testing.T) {
	h, _, _ := newAuthHandler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/auth/login", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNotAcceptable {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestLoginPost_JSONSuccess(t *testing.T) {
	h, acct, _ := newAuthHandler(t)

	payload, _ := json.Marshal(map[string]string{
		"email":    "user@example.com",
		"password": "correct-password",
		"origin":   "/dashboard",
	})
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/login", bytes.NewReader(payload))
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if cookie := w.Result().Cookies(); len(cookie) < 1 || cookie[0].Name != constants.CookieSession || cookie[0].Value == "" {
		t.Fatalf("cookies = %#v", w.Result().Cookies())
	}

	var dto account.User
	if err := json.NewDecoder(w.Body).Decode(&dto); err != nil {
		t.Fatal(err)
	}
	if dto.ID != acct.ID || dto.Email != acct.Email {
		t.Fatalf("dto = %#v", dto)
	}
}

func TestLoginPost_JSONUnauthorized(t *testing.T) {
	h, _, _ := newAuthHandler(t)

	payload, _ := json.Marshal(map[string]string{
		"email":    "user@example.com",
		"password": "wrong",
	})
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/login", bytes.NewReader(payload))
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestLoginPost_HTMLRedirect(t *testing.T) {
	h, _, _ := newAuthHandler(t)

	form := url.Values{}
	form.Set("email", "user@example.com")
	form.Set("password", "correct-password")
	form.Set("origin", "/admin")
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set(constants.HeaderAccept, "text/html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); loc != "/admin" {
		t.Fatalf("Location = %q", loc)
	}
	if cookie := w.Result().Cookies(); len(cookie) < 1 || cookie[0].Name != constants.CookieSession {
		t.Fatalf("cookies = %#v", w.Result().Cookies())
	}
}
