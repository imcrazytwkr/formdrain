package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/routes/health"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

func TestLivez_OK(t *testing.T) {
	t.Parallel()

	h := newHealthHandler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestReadyz_OK(t *testing.T) {
	t.Parallel()

	h := newHealthHandler(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestReadyz_UnavailableWhenDBClosed(t *testing.T) {
	t.Parallel()

	db := testutil.OpenSqlite(t)
	err := db.Close()
	if err != nil {
		t.Fatal(err)
	}

	router := chi.NewRouter()
	health.NewHealthRouter(db).Router(router)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func newHealthHandler(t *testing.T) http.Handler {
	t.Helper()

	db := testutil.OpenSqlite(t)
	router := chi.NewRouter()
	health.NewHealthRouter(db).Router(router)
	return router
}
