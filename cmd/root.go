package cmd

import (
	"os"
	"time"

	"github.com/rendau/dop/adapters/jwk"
	"github.com/rendau/dop/adapters/jwk/jwks"
	dopLoggerZap "github.com/rendau/dop/adapters/logger/zap"
	dopServerHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/dop/dopTools"
	"github.com/rendau/websocket/docs"
	"github.com/rendau/websocket/internal/adapters/server/rest"
	"github.com/rendau/websocket/internal/domain/core"
)

func Execute() {
	var err error

	app := struct {
		lg         *dopLoggerZap.St
		jwk        jwk.Jwk
		core       *core.St
		restApiSrv *dopServerHttps.St
	}{}

	confLoad()

	app.lg = dopLoggerZap.New(conf.LogLevel, conf.Debug)

	if conf.JwkUrl != "" {
		app.jwk, err = jwks.NewByUrl(app.lg, conf.JwkUrl, time.Hour)
		if err != nil {
			app.lg.Fatal(err)
		}
	} else {
		app.lg.Fatal("Jwk-url required")
	}

	app.core = core.New(app.lg, app.jwk)

	docs.SwaggerInfo.Host = conf.SwagHost
	docs.SwaggerInfo.BasePath = conf.SwagBasePath
	docs.SwaggerInfo.Schemes = []string{conf.SwagSchema}
	docs.SwaggerInfo.Title = "Stg service"

	// START

	app.lg.Infow("Starting")

	app.restApiSrv = dopServerHttps.Start(
		conf.HttpListen,
		rest.GetHandler(
			app.lg,
			app.core,
			conf.HttpCors,
		),
		app.lg,
	)

	var exitCode int

	select {
	case <-dopTools.StopSignal():
	case <-app.restApiSrv.Wait():
		exitCode = 1
	}

	// STOP

	app.lg.Infow("Shutting down...")

	if !app.restApiSrv.Shutdown(20 * time.Second) {
		exitCode = 1
	}

	app.lg.Infow("Exit")

	os.Exit(exitCode)
}
