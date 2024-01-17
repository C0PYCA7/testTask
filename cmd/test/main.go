package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"testTask/internal/config"
	"testTask/internal/database/postgres"
	"testTask/internal/handler/user/create"
)

func main() {
	cfg := config.MustLoad()

	log := newLogger()

	log.Info("starting test task", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	database, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to init database ", err)
		os.Exit(1)
	}

	if err := database.Migrate(); err != nil {
		log.Error("failed to create table ", err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/create", create.New(log))

	log.Info("starting server", slog.String("address", cfg.Server.Address))

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}

func newLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
	return log
}
