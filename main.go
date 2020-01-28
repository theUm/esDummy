package main

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/theUm/esDummy/config"
	"github.com/theUm/esDummy/elastic"
	"github.com/theUm/esDummy/srv/healthcheck"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	cfg, err := config.ReadEnv()
	if err != nil {
		panic(errors.Wrap(err, "cant read envs"))
	}
	initLogger(cfg.Log)

	//main context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	var wg = &sync.WaitGroup{}

	esClient, err := elastic.New(ctx, cfg.ElasticConfig)
	if err != nil {
		panic(errors.Wrap(err, "cant connect to elastic"))
	}

	healthCheckService := healthcheck.New(cfg.HealthCheckHTTPPort, esClient.Check)
	healthCheckService.Run(ctx, wg)

	wg.Wait()
	log.Info("terminated")
}

func initLogger(cfg config.LoggerConfig) {
	if cfg.Pretty {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetOutput(os.Stderr)

	switch strings.ToLower(cfg.LogLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Error("Got Interrupt signal")
		stop()
	}()
}
