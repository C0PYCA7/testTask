package main

import (
	"log/slog"
	"os"
	"testTask/internal/config"
)

func main() {
	cfg := config.MustLoad()

	log := newLogger()

	log.Info("starting test task", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	log.Info("config data: ", slog.Any("cfg", cfg))
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
