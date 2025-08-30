package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"rest-api-service/config/storage/sqlite"
	"rest-api-service/internal/config"
	"rest-api-service/internal/config/lib/logger/sl"
	"rest-api-service/internal/http-server/handlers/redirect"
	"rest-api-service/internal/http-server/handlers/save"
	"rest-api-service/internal/http-server/handlers/urldelete"
	"rest-api-service/internal/http-server/middleware/slogpretty"

	st "rest-api-service/config/storage"

	mwLogger "rest-api-service/internal/http-server/middleware/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// Logger
	logger := setupLogger(cfg.Env)
	logger.Info("starting url-shortener", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	// Storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// Router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		router.Post("/url", save.New(logger, storage))
		// TODO: add Delete /url/{uid}
	})

	router.Post("/url", save.New(logger, storage))
	router.Get("/{alias}", redirect.New(logger, storage))
	router.Delete("/url/{alias}", urldelete.New(logger, storage))

	logger.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}
	logger.Error("server stopped")

	alias := "google"
	url := "https://google.com"

	id, err := storage.SaveURL(url, alias)
	if err != nil {
		if errors.Is(err, st.ErrURLExists) { // теперь правильно
			logger.Warn("url already exists", slog.String("alias", alias))
			existingURL, err := storage.GetURL(alias)
			if err != nil {
				logger.Error("failed to get existing URL", sl.Err(err))
			} else {
				logger.Info("existing url found", slog.String("alias", alias), slog.String("url", existingURL))
			}
		} else {
			logger.Error("failed to save url", sl.Err(err))
			os.Exit(1)
		}
	} else {
		logger.Info("saved url", slog.Int64("id", id))
	}

	// Опционально: удалить URL
	deleted, err := storage.DeleteURL(alias)
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
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
