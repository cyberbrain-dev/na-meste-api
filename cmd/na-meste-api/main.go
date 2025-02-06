package main

import (
	"log/slog"
	"os"

	"github.com/cyberbrain-dev/na-meste-api/internal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("Hello")

	// TODO: init database
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
