package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/middleware"
	m "github.com/imcrazytwkr/formdrain/models/http"
	fcr "github.com/imcrazytwkr/formdrain/repositories/form_config"
	frr "github.com/imcrazytwkr/formdrain/repositories/form_response"
	scr "github.com/imcrazytwkr/formdrain/repositories/site_config"
	"github.com/imcrazytwkr/formdrain/routes/form"
	cvs "github.com/imcrazytwkr/formdrain/services/captcha_validation"
	ns "github.com/imcrazytwkr/formdrain/services/notification"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if gin.IsDebugging() {
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

	engine := gin.New()
	engine.Use(middleware.DefaultLogger(), gin.Recovery())
	engine.Use(middleware.RequestId())
	engine.Use(middleware.ResponseFormatParser(m.ContentTypeHTML, m.ContentTypeJSON))

	mongoDb, err := getMongoDb()
	if err != nil {
		log.Fatal().Err(err)
	}

	formConfigRepository := fcr.NewMongoFormConfigRepository(mongoDb)
	formResponseRepository := frr.NewMongoFormResponseRepository(mongoDb)

	siteConfigRepository := scr.NewMongoSiteConfigRepository(mongoDb)

	captchaValidationService := cvs.NewHttpCaptchaValidationService(http.DefaultClient, &log.Logger)
	notificationService := ns.NewHttpNotificationService(http.DefaultClient)

	form.NewFormRouter(
		formConfigRepository,
		siteConfigRepository,
		formResponseRepository,
		captchaValidationService,
		notificationService,
	).Register(engine.Group("/form"))

	engine.Run(fmt.Sprintf("%s:%s", listenHost, listenPort))
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
