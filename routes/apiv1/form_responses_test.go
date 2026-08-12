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
