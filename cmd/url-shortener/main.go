package main

import (
	"fmt"
	"log/slog"
	"os"
	"rest-api-service/config/storage/sqlite"
	"rest-api-service/internal/config"
	"rest-api-service/internal/config/lib/logger/sl"

	"github.com/go-chi/chi"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	logger := setupLogger(cfg.Env)
	logger.Info("starting url-shortener", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)

	}

	router := chi.NewRouter

	id, err := storage.SaveURL("https://google.com", "google")
	if err != nil {
		logger.Error("failed to save url", sl.Err(err))
		os.Exit(1)
	}

	logger.Info("saved url", slog.Int64("id", id))

	id, err = storage.SaveURL("https://google.com", "google")
	if err != nil {
		logger.Error("failed to save url", sl.Err(err))
		os.Exit(1)
	}

	deleted, err := storage.DeleteURL("google")
	if err != nil {
		logger.Error("failed to delete url", sl.Err(err))
	} else {
		logger.Info("deleted url", slog.Int64("rows", deleted))
	}

	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return log
}
