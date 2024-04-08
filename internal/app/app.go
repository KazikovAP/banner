package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"banner/internal/config"
	"banner/internal/handler"
	api "banner/internal/http"
	"banner/internal/repository"
	"banner/internal/repository/postgres"
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
	log.Info("starting banner-server", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Setup connect to database
	db, err := setupConnectToPostgres(cfg.Postgres)
	if err != nil {
		log.Error("failed to connect Postgres: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	service := handler.NewBannerService(db)
	handler := api.NewHandler(service)

	httpServer := setupServer(cfg.Server, handler)

	log.Debug("starting HTTP server on %s", httpServer.Addr)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Error("failed listen and serve: %v", err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	log.Debug("shutting down server")
	if err := httpServer.Shutdown(context.Background()); err != nil {
		return err
	}
	if err := db.Close(context.Background()); err != nil {
		return err
	}

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

func setupConnectToPostgres(cfg config.PostgresConfig) (repository.BannerStorage, error) {
	log.Println("setup storage")

	storage := fmt.Sprintf("host=%s port=%s user=%s password=%s db=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database)

	return postgres.NewPostgres(storage)
}

func setupServer(cfg config.ServerConfig, handler *api.Handler) *http.Server {
	log.Println("setup HTTP server")

	router := api.Router(handler)

	return &http.Server{
		Addr:    fmt.Sprintf("[::]:%s", cfg.Port),
		Handler: router,
	}
}
