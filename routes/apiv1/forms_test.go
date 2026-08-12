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

func seedOwnedForms(t *testing.T, h *apiV1Harness) {
	t.Helper()
	h.seedSite(1, "c.example")
	h.seedSite(2, "a.example")
	h.seedSite(3, "b.example")
	testutil.SeedForm(t, h.db, 11, 1)
	testutil.SeedForm(t, h.db, 12, 1)
	testutil.SeedForm(t, h.db, 21, 2)
	testutil.SeedForm(t, h.db, 31, 3)
}

func TestListForms_ByID(t *testing.T) {
	h := newApiV1Harness(t)
	seedOwnedForms(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?limit=2", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.FormList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 2 || page.Items[0].Id != 11 || page.Items[1].Id != 12 {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.NextCursor == nil || len(*page.NextCursor) < 1 {
		t.Fatal("expected next_cursor")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?limit=2&cursor="+*page.NextCursor, nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page2 api.FormList
	if err := json.NewDecoder(w.Body).Decode(&page2); err != nil {
		t.Fatal(err)
	}
	if len(page2.Items) != 2 || page2.Items[0].Id != 21 || page2.Items[1].Id != 31 {
		t.Fatalf("page2 = %#v", page2.Items)
	}
	if page2.NextCursor != nil {
		t.Fatalf("unexpected next_cursor %q", *page2.NextCursor)
	}
}

func TestListForms_BySiteID(t *testing.T) {
	h := newApiV1Harness(t)
	seedOwnedForms(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?sort=site_id&limit=2", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.FormList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 2 || page.Items[0].Id != 11 || page.Items[1].Id != 12 {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.Items[0].SiteId != 1 || page.Items[1].SiteId != 1 {
		t.Fatalf("site_ids = %#v", page.Items)
	}
	if page.NextCursor == nil {
		t.Fatal("expected next_cursor")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?sort=site_id&limit=2&cursor="+*page.NextCursor, nil)
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
	if len(page.Items) != 2 || page.Items[0].Id != 21 || page.Items[1].Id != 31 {
		t.Fatalf("page2 = %#v", page.Items)
	}
}

func TestListForms_ByHostname(t *testing.T) {
	h := newApiV1Harness(t)
	seedOwnedForms(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?sort=hostname&limit=2", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.FormList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	// a.example (21), b.example (31), then c.example (11, 12)
	if len(page.Items) != 2 || page.Items[0].Id != 21 || page.Items[1].Id != 31 {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.NextCursor == nil {
		t.Fatal("expected next_cursor")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?sort=hostname&limit=2&cursor="+*page.NextCursor, nil)
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
	if len(page.Items) != 2 || page.Items[0].Id != 11 || page.Items[1].Id != 12 {
		t.Fatalf("page2 = %#v", page.Items)
	}
}

func TestListForms_FilterBySiteID(t *testing.T) {
	h := newApiV1Harness(t)
	seedOwnedForms(t, h)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?site_id=1&limit=1", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.FormList
	if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 || page.Items[0].Id != 11 || page.Items[0].SiteId != 1 {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.NextCursor == nil {
		t.Fatal("expected next_cursor")
	}

	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?site_id=1&limit=1&cursor="+*page.NextCursor, nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w = httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page2 api.FormList
	if err := json.NewDecoder(w.Body).Decode(&page2); err != nil {
		t.Fatal(err)
	}
	if len(page2.Items) != 1 || page2.Items[0].Id != 12 {
		t.Fatalf("page2 = %#v", page2.Items)
	}
	if page2.NextCursor != nil {
		t.Fatalf("unexpected next_cursor %q", *page2.NextCursor)
	}
}

func TestListForms_FilterBySiteID_ForeignOrMissing(t *testing.T) {
	h := newApiV1Harness(t)
	seedOwnedForms(t, h)
	testutil.SeedSite(t, h.db, 9, "other.example")
	testutil.SeedForm(t, h.db, 91, 9)

	for _, siteID := range []string{"9", "999"} {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?site_id="+siteID, nil)
		req.Header.Set(constants.HeaderAccept, "application/json")
		req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
		w := httptest.NewRecorder()
		h.handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("site_id=%s status = %d body=%s", siteID, w.Code, w.Body.String())
		}
		var page api.FormList
		if err := json.NewDecoder(w.Body).Decode(&page); err != nil {
			t.Fatal(err)
		}
		if len(page.Items) != 0 {
			t.Fatalf("site_id=%s items = %#v", siteID, page.Items)
		}
	}
}

func TestListForms_Empty(t *testing.T) {
	h := newApiV1Harness(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	var page api.FormList
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

func TestListForms_BadSort(t *testing.T) {
	h := newApiV1Harness(t)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms?sort=nope", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestListForms_BadCursor(t *testing.T) {
	h := newApiV1Harness(t)

	cases := []string{
		"/api/v1/forms?cursor=not-valid",
		"/api/v1/forms?sort=site_id&cursor=not-valid",
		"/api/v1/forms?sort=hostname&cursor=not-valid",
		"/api/v1/forms?site_id=1&cursor=not-valid",
	}
	h.seedSite(1, "owned.example")

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

func TestListForms_Unauthorized(t *testing.T) {
	h := newApiV1Harness(t)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/forms", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}
