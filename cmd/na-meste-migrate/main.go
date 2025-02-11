package main

import (
	"log/slog"
	"os"

	"github.com/cyberbrain-dev/na-meste-api/internal/config"
	"github.com/cyberbrain-dev/na-meste-api/internal/database"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// loading th config
	cfg := config.MustLoad()

	// launching the slogger
	logger := setupLogger(cfg.Env)

	// some info
	logger.Info("Launching the migration utility...")

	// connecting to the db
	logger.Info("Connecting to Postgres database...")
	db, err := database.ConnectPostgres(cfg.PostgresConnection)
	if err != nil {
		// logging the error
		logger.Error(
			"Connection was not successful",
			slog.Any("err", err),
		)
		os.Exit(1)
	}
	logger.Info("Successfuly connected to Postgres database")

	// migrating...
	logger.Info("Starting migrating entities...")
	if err := database.MigrateEntities(db); err != nil {
		logger.Error("Failed to migrate entities")
	} else {
		logger.Info("Migration has been successful")
	}

	// disconnecting...
	err = database.DisconnectPostgres(db)
	if err != nil {
		// logging the error
		logger.Error(
			"Unable to close connection to Postgres database",
			slog.Any("err", err),
		)
	} else {
		logger.Info("Successfuly disconnected Postgres database")
	}
}

// Sets up a slog logger
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
