package main

import (
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/imcrazytwkr/formdrain/middleware"
	m "github.com/imcrazytwkr/formdrain/models/http"
	fcr "github.com/imcrazytwkr/formdrain/repositories/form_config"
	frr "github.com/imcrazytwkr/formdrain/repositories/form_response"
	scr "github.com/imcrazytwkr/formdrain/repositories/site_config"
	"github.com/imcrazytwkr/formdrain/routes/form"
	cvs "github.com/imcrazytwkr/formdrain/services/captcha_validation"
	ns "github.com/imcrazytwkr/formdrain/services/notification"
	"github.com/imcrazytwkr/formdrain/utils/httpclient"
	"github.com/imcrazytwkr/formdrain/utils/httpserver"
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

	err = httpserver.LoadTemplates()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load templates")
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

	httpClient := httpclient.DefaultClient()
	captchaValidationService := cvs.NewHttpCaptchaValidationService(httpClient, &log.Logger)
	notificationService := ns.NewHttpNotificationService(httpClient)

	router.Route("/form", func(r chi.Router) {
		form.NewFormRouter(
			formConfigRepository,
			siteConfigRepository,
			formResponseRepository,
			captchaValidationService,
			notificationService,
		).Register(r)
	})

	addr := net.JoinHostPort(listenHost, listenPort)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal().Err(err).Msg("server stopped")
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

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
