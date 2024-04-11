package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"banner/internal/config"
	middlewares "banner/internal/lib/api/middlewares"
	jwt "banner/internal/lib/auth/jwt"
	logerr "banner/internal/lib/logger/logerr"
	"banner/internal/repo"
	"banner/internal/repository/postgres"
	"banner/internal/server/handlers/banners"
	"banner/internal/server/handlers/features"
	"banner/internal/server/handlers/tags"
	login "banner/internal/server/handlers/users/login"
	users "banner/internal/server/handlers/users/user"

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
		log.Error("Failed to connect Postgres", logerr.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Error("Failed to ping Postgres", logerr.Err(err))
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
	router.Post("/features", features.NewFeature(log, ftr))

	tg := repo.NewTag(db.DB, log)
	router.Post("/tags", tags.NewTag(log, tg))

	us := repo.NewUser(db.DB, log)
	router.Post("/users", users.NewUser(log, us))

	jwt := jwt.NewJWTSecret(cfg.Jwt.Secret, log)
	router.Post("/login", login.Login(log, us, jwt))

	router.With(func(next http.Handler) http.Handler {
		return middlewares.TokenAuthMiddleware(jwt, next)
	}).Post("/tags", tags.NewTag(log, tg))

	br := repo.NewBanner(db.DB, log)
	btr := repo.NewBannerTag(db.DB, log)

	router.With(func(next http.Handler) http.Handler {
		return middlewares.TokenAuthAndRoleMiddleware(jwt, next)
	}).Post("/banners", banners.NewBanner(log, br, btr))

	// Server
	log.Info("Starting server at", slog.String(cfg.Server.Host, cfg.Server.Port))
	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start server", logerr.Err(err))
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
