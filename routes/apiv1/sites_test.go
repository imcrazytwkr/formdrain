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

func TestGetSiteConfig_OK(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(7, "owned.example")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites/7", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}

	var got api.Site
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Id != 7 || got.Hostname != "owned.example" || got.OwnerId != h.account.ID {
		t.Fatalf("got = %#v", got)
	}
}

func TestGetSiteConfig_NotFound(t *testing.T) {
	h := newApiV1Harness(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites/999", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetSiteConfig_OtherOwner(t *testing.T) {
	h := newApiV1Harness(t)
	testutil.SeedSite(t, h.db, 8, "other.example")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites/8", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetSiteConfig_Unauthorized(t *testing.T) {
	h := newApiV1Harness(t)
	h.seedSite(7, "owned.example")

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites/7", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)

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

func TestListSites_Empty(t *testing.T) {
	h := newApiV1Harness(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites", nil)
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
	if len(page.Items) != 0 {
		t.Fatalf("items = %#v", page.Items)
	}
	if page.NextCursor != nil {
		t.Fatalf("unexpected next_cursor %q", *page.NextCursor)
	}
}

func TestListSites_BadSort(t *testing.T) {
	h := newApiV1Harness(t)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites?sort=nope", nil)
	req.Header.Set(constants.HeaderAccept, "application/json")
	req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: h.sessionID})
	w := httptest.NewRecorder()
	h.handler.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
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

func TestListSites_BadHostnameCursor(t *testing.T) {
	h := newApiV1Harness(t)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/sites?sort=hostname&cursor=not-valid", nil)
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
