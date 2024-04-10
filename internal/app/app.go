package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"banner/internal/config"
	"banner/internal/lib/logger"
	"banner/internal/repo"
	"banner/internal/repository/postgres"
	"banner/internal/server/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run() error {
	// Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	fmt.Println(cfg.Env)

	// Logger
	log := setupLogger(cfg.Env)
	log.Info("Starting banner-server", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	// Setup connect to database
	db, err := setupConnectToPostgres(cfg, log)
	if err != nil {
		log.Error("Failed to connect Postgres", logger.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Error("Failed to ping Postgres", logger.Err(err))
		os.Exit(1)
	} else {
		log.Info("Connection to Postgres DB successfully")
	}
	log.Info("Application started...", slog.String("env", cfg.Env))

	// Router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	ftr := repo.NewFeature(db.DB, log)
	router.Post("/features", handlers.NewFeature(log, ftr))

	log.Info("Starting server at", slog.String(cfg.Server.Host, cfg.Server.Port))
	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start server", logger.Err(err))
	}

	// Server TODO

	return nil
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupConnectToPostgres(cfg *config.Config, log *slog.Logger) (*postgres.Postgres, error) {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database)

	db, err := postgres.NewPostgres(context.Background(), connection, log)

	return db, err
}
