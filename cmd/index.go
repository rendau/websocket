package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rendau/websocket/internal/adapters/api/rest"
	"github.com/rendau/websocket/internal/adapters/logger/zap"
	"github.com/rendau/websocket/internal/domain/core"
	"github.com/spf13/viper"
)

func Execute() {
	loadConf()

	debug := viper.GetBool("debug")

	lg, err := zap.New(viper.GetString("log_level"), debug, false)
	if err != nil {
		log.Fatal(err)
	}
	defer lg.Sync()

	cr, err := core.New(lg)
	if err != nil {
		lg.Fatal(err)
	}

	httpApi := rest.New(lg, viper.GetString("http_listen"), viper.GetString("usr_auth_url"), cr)

	lg.Infow("Starting", "http_listen", viper.GetString("http_listen"))

	httpApi.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var exitCode int

	select {
	case <-stop:
	case <-httpApi.Wait():
		exitCode = 1
	}

	lg.Infow("Shutting down...")

	ctx, ctxCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer ctxCancel()

	err = httpApi.Shutdown(ctx)
	if err != nil {
		lg.Errorw("Fail to shutdown http-api", err)
		exitCode = 1
	}

	lg.Infow("Waiting running jobs...")

	os.Exit(exitCode)
}

func loadConf() {
	viper.SetDefault("debug", "false")
	viper.SetDefault("http_listen", ":80")
	viper.SetDefault("log_level", "warn")

	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath == "" {
		confFilePath = "conf.yml"
	}
	viper.SetConfigFile(confFilePath)
	_ = viper.ReadInConfig()

	// env vars are in priority
	viper.AutomaticEnv()
}
