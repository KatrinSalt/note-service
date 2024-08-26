package main

import (
	"fmt"
	"os"

	"github.com/KatrinSalt/notes-service/config"
	"github.com/KatrinSalt/notes-service/log"
	"github.com/KatrinSalt/notes-service/server"
)

func main() {
	logger, err := setLogger()
	if err != nil {
		logger.Error("Server error.", "error", err)
		os.Exit(1)
	}

	if err := run(logger); err != nil {
		logger.Error("Server error.", "error", err)
		os.Exit(1)
	}
}

func setLogger() (*log.Logger, error) {
	logLevel, set := os.LookupEnv("SERVER_LOG_LEVEL")

	if set {
		logger, err := log.NewWithSetLevel(logLevel)
		if err != nil {
			return nil, fmt.Errorf("could not load configuration: %w", err)
		}
		return logger, nil
	} else {
		return log.New(), nil
	}

}
func run(log *log.Logger) error {
	log.Info("Starting the note service.")
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("could not load configuration: %w", err)
	}

	log.Debug("Loaded configuration.", "config", cfg)

	services, err := config.SetupServices(cfg.Services)
	if err != nil {
		return fmt.Errorf("could not setup services: %w", err)
	}

	srv, err := server.New(
		services.Note,
		server.WithAddress(cfg.Server.Host+":"+cfg.Server.Port),
		server.WithLogger(log),
		server.WithLogger(log),
	)
	if err != nil {
		return fmt.Errorf("could not create server: %w", err)
	}

	if err := srv.Start(); err != nil {
		return fmt.Errorf("could not start server: %w", err)

	}
	log.Info("Note service stopped.")
	return nil
}
