package main

import (
	"context"
	"errors"
	"github.com/elgntt/effective-mobile/internal/api"
	"github.com/elgntt/effective-mobile/internal/config"
	db2 "github.com/elgntt/effective-mobile/internal/pkg/db"
	"github.com/elgntt/effective-mobile/internal/repository"
	"github.com/elgntt/effective-mobile/internal/service"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var log = logrus.New()

func main() {
	log.SetOutput(os.Stdout)
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ServerConfig.ServerMode == "debug" {
		log.SetLevel(logrus.DebugLevel)
	}

	log.Debugln("Connecting to the database...")
	db, err := db2.OpenDB(context.Background(), cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	log.Debugln("Up migrations...")
	err = db2.Migrate(db, cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}

	r := api.New(
		service.New(
			repository.New(db),
			cfg.APIConfig,
			log,
		),
		log,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerConfig.ServerPort,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Infoln("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Infoln("Server exiting")
}
