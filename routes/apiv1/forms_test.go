package apiv1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestGetFormConfig_OK(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(1, "owned.example")
	testutil.SeedForm(t, h.db, 10, 1)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var got api.FormConfig
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Id != 10 || got.SiteId != 1 || got.CaptchaType != api.Hcaptcha {
		t.Fatalf("got = %#v", got)
	}
}

func TestGetFormConfig_NotFound(t *testing.T) {
	h := newApiV1Harness(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/999", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetFormConfig_OtherOwner(t *testing.T) {
	h := newApiV1Harness(t)
	testutil.SeedSite(t, h.db, 8, "other.example")
	testutil.SeedForm(t, h.db, 80, 8)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/80", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetFormConfig_Unauthorized(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(1, "owned.example")
	testutil.SeedForm(t, h.db, 10, 1)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms/10", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}
