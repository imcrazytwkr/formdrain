package apiv1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/models/account"
)

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
