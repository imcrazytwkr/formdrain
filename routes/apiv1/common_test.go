package apiv1_test

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/imcrazytwkr/formdrain/models/account"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/models/session"
	ar "github.com/imcrazytwkr/formdrain/repositories/account"
	fcr "github.com/imcrazytwkr/formdrain/repositories/form_config"
	sr "github.com/imcrazytwkr/formdrain/repositories/session"
	scr "github.com/imcrazytwkr/formdrain/repositories/site_config"
	"github.com/imcrazytwkr/formdrain/routes/apiv1"
	as "github.com/imcrazytwkr/formdrain/services/account"
	"github.com/imcrazytwkr/formdrain/utils/testutil"
)

type apiV1Harness struct {
	handler   http.Handler
	account   *account.Account
	sessionID string
	db        *sql.DB
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
		db:        db,
		seedSite: func(siteID int64, hostname string) {
			testutil.SeedSiteForOwner(t, db, acct.ID, siteID, hostname)
		},
	}
}
