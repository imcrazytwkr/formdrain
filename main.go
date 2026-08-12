package main

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/middleware"
	m "github.com/imcrazytwkr/formdrain/models/http"
	ar "github.com/imcrazytwkr/formdrain/repositories/account"
	fcr "github.com/imcrazytwkr/formdrain/repositories/form_config"
	frr "github.com/imcrazytwkr/formdrain/repositories/form_response"
	sr "github.com/imcrazytwkr/formdrain/repositories/session"
	scr "github.com/imcrazytwkr/formdrain/repositories/site_config"
	"github.com/imcrazytwkr/formdrain/routes/apiv1"
	"github.com/imcrazytwkr/formdrain/routes/auth"
	"github.com/imcrazytwkr/formdrain/routes/form"
	as "github.com/imcrazytwkr/formdrain/services/account"
	cvs "github.com/imcrazytwkr/formdrain/services/captcha_validation"
	ns "github.com/imcrazytwkr/formdrain/services/notification"
	"github.com/imcrazytwkr/formdrain/utils/httpclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "modernc.org/sqlite"
)

func main() {
	if os.Getenv("LOG_MODE") == "release" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	listenHost, err := getHost()
	if err != nil {
		log.Fatal().Err(err)
	}

	listenPort, err := getPort()
	if err != nil {
		log.Fatal().Err(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.DefaultLogger())
	router.Use(middleware.RequestId())
	router.Use(middleware.ResponseFormatParser(m.ContentTypeHTML, m.ContentTypeJSON))
	router.Use(middleware.Recoverer())

	sqliteDB, err := openSqlite()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open sqlite database")
	}
	defer sqliteDB.Close()

	formConfigRepository := fcr.NewSqliteFormConfigRepository(sqliteDB)
	formResponseRepository := frr.NewSqliteFormResponseRepository(sqliteDB)
	siteConfigRepository := scr.NewSqliteSiteConfigRepository(sqliteDB)
	accountRepository := ar.NewSqliteAccountRepository(sqliteDB)
	sessionRepository := sr.NewMemorySessionRepository()

	httpClient := httpclient.DefaultClient()
	captchaValidationService := cvs.NewHttpCaptchaValidationService(httpClient, &log.Logger)
	accountService := as.NewService(accountRepository)

	brevoAPIKey, err := getBrevoAPIKey()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load Brevo API key")
	}

	notificationService := ns.NewHttpNotificationService(httpClient, brevoAPIKey)

	router.Route("/auth",
		auth.NewAuthRouter(
			sessionRepository,
			accountService,
		).Router)
	router.Route("/api/v1",
		apiv1.NewApiV1Router(
			sessionRepository,
			accountRepository,
			siteConfigRepository,
			formConfigRepository,
			formResponseRepository,
		).Router)
	router.Route("/form",
		form.NewFormRouter(
			formConfigRepository,
			siteConfigRepository,
			formResponseRepository,
			captchaValidationService,
			notificationService,
		).Router,
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    net.JoinHostPort(listenHost, listenPort),
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server stopped")
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), terminationGracePeriod)
	defer cancel()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatal().Err(err).Msg("server shutdown")
	}
}

func openSqlite() (*sql.DB, error) {
	dbURL, err := getDBURL()
	if err != nil {
		return nil, err
	}

	path, err := sqliteFilePath(dbURL)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
