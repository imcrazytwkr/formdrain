package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

	mongoDb, err := getMongoDb()
	if err != nil {
		log.Fatal().Err(err)
	}

	formConfigRepository := fcr.NewMongoFormConfigRepository(mongoDb)
	formResponseRepository := frr.NewMongoFormResponseRepository(mongoDb)

	siteConfigRepository := scr.NewMongoSiteConfigRepository(mongoDb)

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

func getMongoDb(opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	connstring, err := getConnString()
	if err != nil {
		return nil, err
	}

	database := strings.Trim(connstring.Path, "/")
	if len(database) < 1 {
		database = "formdrain"
	}

	mongoOptions := options.Client().ApplyURI(connstring.String())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		return nil, err
	}

	return mongoClient.Database(database, opts...), nil
}
