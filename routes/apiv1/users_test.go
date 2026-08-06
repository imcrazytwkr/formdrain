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
	sr "github.com/imcrazytwkr/formdrain/repositories/session"
	"github.com/imcrazytwkr/formdrain/routes/apiv1"
	as "github.com/imcrazytwkr/formdrain/services/account"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func newApiV1Handler(t *testing.T) (http.Handler, *account.Account, string) {
	t.Helper()

	db := testutil.OpenSqlite(t)
	accounts := ar.NewSqliteAccountRepository(db)
	sessions := sr.NewMemorySessionRepository()

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
	router.Route("/api/v1", apiv1.NewApiV1Router(sessions, accounts).Router)

	return router, acct, sess.ID
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
