package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
)

func TestLogoutPost_JSONSuccess(t *testing.T) {
	h, _, sessions := newAuthHandler(t)
	sessionID := loginSession(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/logout", strings.NewReader(`{}`))
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: sessionID})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	cleared := false
	for _, c := range w.Result().Cookies() {
		if c.Name == constants.CookieSession {
			cleared = true
			if c.Value != "" || c.MaxAge >= 0 {
				t.Fatalf("cookie = %#v", c)
			}
		}
	}
	if !cleared {
		t.Fatal("expected cleared session cookie")
	}

	stored, err := sessions.GetByID(t.Context(), sessionID)
	if err != nil || stored != nil {
		t.Fatalf("stored=%#v err=%v", stored, err)
	}
}

func TestLogoutPost_JSONWithoutCookie(t *testing.T) {
	h, _, _ := newAuthHandler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/logout", strings.NewReader(`{}`))
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestLogoutPost_HTMLRedirect(t *testing.T) {
	h, _, sessions := newAuthHandler(t)
	sessionID := loginSession(t, h)

	form := url.Values{}
	form.Set("origin", "/login")
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/logout", strings.NewReader(form.Encode()))
	req.Header.Set(constants.HeaderAccept, "text/html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: sessionID})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); loc != "/login" {
		t.Fatalf("Location = %q", loc)
	}

	stored, err := sessions.GetByID(t.Context(), sessionID)
	if err != nil || stored != nil {
		t.Fatalf("stored=%#v err=%v", stored, err)
	}
}

func TestLogoutPost_HTMLInvalidOriginBecomesRoot(t *testing.T) {
	h, _, _ := newAuthHandler(t)

	payload, _ := json.Marshal(map[string]string{"origin": "https://evil.example/"})
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/logout", bytes.NewReader(payload))
	req.Header.Set(constants.HeaderAccept, "text/html")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); loc != "/" {
		t.Fatalf("Location = %q", loc)
	}
}
