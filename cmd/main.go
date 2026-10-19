package main

import (
	"context"
	"os"
	"syscall"

	"singel-sondag/config"
	"singel-sondag/http"
	"singel-sondag/repo"

	"github.com/jesperkha/notifier"
)

func main() {
	config := config.Load()
	notif := notifier.New()

	logger := repo.NewLogger(config.Environment)
	database, err := repo.New(logger, config.DatabaseUrl)
	if err != nil {
		logger.Error(context.Background(), "failed to connect to database", err)
		os.Exit(1)
	}

	if err := database.Migrate("./migrations"); err != nil {
		logger.Error(context.Background(), "failed to migrate database", err)
		os.Exit(1)
	}

	deps := http.Dependencies{
		Logger:       logger,
		Notif:        notif,
		Config:       config,
		RegisterRepo: database,
	}

	go http.Run(deps)

	notif.NotifyOnSignal(os.Interrupt, syscall.SIGTERM)
	logger.Info(context.Background(), "shutdown complete")
}
