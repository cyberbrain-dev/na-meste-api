package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyberbrain-dev/na-meste-api/internal/config"
	"github.com/cyberbrain-dev/na-meste-api/internal/database"
	"github.com/cyberbrain-dev/na-meste-api/internal/database/repositories"
	"github.com/cyberbrain-dev/na-meste-api/internal/server/endpoints"
	myMw "github.com/cyberbrain-dev/na-meste-api/internal/server/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	logger.Info("connecting to Postgres database...")

	// connecting to the db
	db, err := database.ConnectPostgres(cfg.PostgresConnection)
	if err != nil {
		// logging the error
		logger.Error(
			"connection was not successful",
			slog.Any("err", err),
		)
		os.Exit(1)
	}

	rc := repositories.NewColleges(db)
	ru := repositories.NewUsers(db)
	ra := repositories.NewAttendances(db)

	logger.Info("successfuly connected to Postgres database")

	// initializing a router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	// ! settin' up the routes

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Все на месте!"))
	})

	router.Post("/colleges/", endpoints.CreateCollege(logger, rc))
	router.Post("/users/", endpoints.Register(logger, ru))
	router.Post("/login/", endpoints.Login(logger, ru))

	// registring the attendance creation endpoint and setting a middleware
	router.Post("/attendances/", myMw.CheckRole(
		logger,
		"scanner",
		endpoints.CreateAttendance(
			logger, ra,
		),
	),
	)
	// !

	logger.Info(
		"launching the server...",
		slog.String("address", cfg.HTTPServer.Address),
	)

	// creating a channel that is going to read interrupt signals
	done := make(chan os.Signal, 1)
	// making the signals be written in the channel
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// creating an HTTP server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// running the server in separate gorutine
	// so the code that is under this goroutine can be executed
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Warn("server not running", slog.Any("err", err))
		}
	}()

	logger.Info("server started")

	// waiting for signals
	// (and blocking the execution in the goroutine of main fuction)
	<-done
	// once there is a signal the program gracefully shutdowns
	fmt.Println()
	logger.Info("stopping server...")

	// creating a context for shutting down
	// the server with 10s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// gracefully shutting down the server...
	// if it hasn't closed all the connections
	// during the timeout it is shutdowned by force
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", slog.Any("err", err))

		return
	}

	// disconnecting the database
	err = database.DisconnectPostgres(db)
	if err != nil {
		// logging the error
		logger.Error(
			"unable to close connection to Postgres database",
			slog.Any("err", err),
		)
	} else {
		logger.Info("successfuly disconnected Postgres database")
	}

	// final log
	logger.Info("server stopped...")
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
