package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/imcrazytwkr/formdrain/models/session"
	sr "github.com/imcrazytwkr/formdrain/repositories/session"
)

func TestAuthenticated_MissingCookie(t *testing.T) {
	sessions := sr.NewMemorySessionRepository()
	called := false
	h := middleware.Authenticated(sessions)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
	if called {
		t.Fatal("next should not be called")
	}
}

func TestAuthenticated_UnknownSession(t *testing.T) {
	sessions := sr.NewMemorySessionRepository()
	called := false
	h := middleware.Authenticated(sessions)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: "missing"})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
	if called {
		t.Fatal("next should not be called")
	}
}

func TestAuthenticated_OK(t *testing.T) {
	sessions := sr.NewMemorySessionRepository()
	sess := &session.Session{
		AccountID: 42,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	if err := sessions.Create(t.Context(), sess); err != nil {
		t.Fatal(err)
	}

	var gotID string
	h := middleware.Authenticated(sessions)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, ok := middleware.SessionFromContext(r.Context())
		if !ok {
			t.Fatal("session missing from context")
		}
		gotID = s.ID
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: sess.ID})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
	if gotID != sess.ID {
		t.Fatalf("session id = %q, want %q", gotID, sess.ID)
	}
}
