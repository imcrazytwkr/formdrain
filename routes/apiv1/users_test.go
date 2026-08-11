package apiv1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/imcrazytwkr/formdrain/models/account"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/models/session"
	ar "github.com/imcrazytwkr/formdrain/repositories/account"
	fcr "github.com/imcrazytwkr/formdrain/repositories/form_config"
	sr "github.com/imcrazytwkr/formdrain/repositories/session"
	scr "github.com/imcrazytwkr/formdrain/repositories/site_config"
	"github.com/imcrazytwkr/formdrain/routes/apiv1"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	as "github.com/imcrazytwkr/formdrain/services/account"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

type apiV1Harness struct {
	handler   http.Handler
	account   *account.Account
	sessionID string
	seedSite  func(siteID int64, hostname string)
}

func newApiV1Handler(t *testing.T) (http.Handler, *account.Account, string) {
	t.Helper()
	h := newApiV1Harness(t)
	return h.handler, h.account, h.sessionID
}

func newApiV1Harness(t *testing.T) *apiV1Harness {
	t.Helper()

	db := testutil.OpenSqlite(t)
	accounts := ar.NewSqliteAccountRepository(db)
	sessions := sr.NewMemorySessionRepository()
	sites := scr.NewSqliteSiteConfigRepository(db)
	forms := fcr.NewSqliteFormConfigRepository(db)

	hash, err := as.HashPassword("correct-password")
	if err != nil {
		t.Fatal(err)
	}

	acct := &account.Account{Email: "user@example.com", PasswordHash: hash}
	if err := accounts.Create(t.Context(), acct); err != nil {
		t.Fatal(err)
	}

	sess := &session.Session{
		AccountID: acct.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	if err := sessions.Create(t.Context(), sess); err != nil {
		t.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.ResponseFormatParser(m.ContentTypeHTML, m.ContentTypeJSON))
	router.Route("/api/v1", apiv1.NewApiV1Router(sessions, accounts, sites, forms).Router)

	return &apiV1Harness{
		handler:   router,
		account:   acct,
		sessionID: sess.ID,
		seedSite: func(siteID int64, hostname string) {
			testutil.SeedSiteForOwner(t, db, acct.ID, siteID, hostname)
		},
	}
}

func TestGetCurrentUser_OK(t *testing.T) {
	h, acct, sessionID := newApiV1Handler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/users/current", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: sessionID})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var got account.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.ID != acct.ID || got.Email != acct.Email {
		t.Fatalf("got = %#v want id=%d email=%s", got, acct.ID, acct.Email)
	}
}

func TestGetCurrentUser_Unauthorized(t *testing.T) {
	h, _, _ := newApiV1Handler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/users/current", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestListSites_ByID(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(1, "c.example")
	h.seedSite(2, "a.example")
	h.seedSite(3, "b.example")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites?limit=2", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.SiteList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 2 || page.Items[0].Id != 1 || page.Items[1].Id != 2 {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.NextCursor == nil || len(*page.NextCursor) < 1 {
		t.Fatal("expected next_cursor")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites?limit=2&cursor="+*page.NextCursor, nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page2 api.SiteList
	if err := json.NewDecoder(w.Body).Decode(&page2); err != nil {
		t.Fatal(err)
	}
	if len(page2.Items) != 1 || page2.Items[0].Id != 3 {
		t.Fatalf("page2 = %#v", page2.Items)
	}
	if page2.NextCursor != nil {
		t.Fatalf("unexpected next_cursor %q", *page2.NextCursor)
	}
}

func TestListSites_ByHostname(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(1, "c.example")
	h.seedSite(2, "a.example")
	h.seedSite(3, "b.example")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites?sort=hostname&limit=2", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.SiteList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 2 || page.Items[0].Hostname != "a.example" || page.Items[1].Hostname != "b.example" {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.NextCursor == nil {
		t.Fatal("expected next_cursor")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites?sort=hostname&limit=2&cursor="+*page.NextCursor, nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 || page.Items[0].Hostname != "c.example" {
		t.Fatalf("page2 = %#v", page.Items)
	}
}

func TestListSites_BadCursor(t *testing.T) {
	h := newApiV1Harness(t)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites?cursor=not-valid", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestListSites_Unauthorized(t *testing.T) {
	h := newApiV1Harness(t)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}
