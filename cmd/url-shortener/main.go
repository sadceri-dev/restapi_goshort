package main

import (
	"LINKSHORTENER/internal/config"
	"LINKSHORTENER/internal/lib/logger/sl"
	"LINKSHORTENER/internal/storage"
	"LINKSHORTENER/internal/storage/sqlite"
	"log/slog"
	"os"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

router :=chi.NewRouter()

router.Use(middleware.RequestID)
router.Use(middleware.Logger)
router.Use(middleware.Recoverer)
router.Use(middleware.URLFormat)


func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("Init server", slog.String("address", cfg.Address))
	log.Debug("Logger debug mode enable")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
