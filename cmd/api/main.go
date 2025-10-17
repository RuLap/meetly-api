package main

import (
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/handlers"
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/repository"
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/services"
	"github.com/RuLap/meetly-api/meetly/internal/pkg/config"
	"github.com/RuLap/meetly-api/meetly/internal/pkg/jwt_helper"
	"github.com/RuLap/meetly-api/meetly/internal/pkg/logger"
	postgres "github.com/RuLap/meetly-api/meetly/internal/pkg/storage"
	validation "github.com/RuLap/meetly-api/meetly/internal/pkg/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.New(logger.Config{
		Level: cfg.Env,
	})

	validation.Init()

	storage, err := postgres.InitDB(cfg.PostgresConnString)
	if err != nil {

	}

	jwtHelper, err := jwt_helper.NewJwtHelper(cfg.JWT.Secret)
	if err != nil {

	}

	authRepo := repository.NewAuthRepository(storage.Database())
	logger.Info("Init repos successfully")

	googleConfig := &services.GoogleOAuthConfig{
		ClientID:     cfg.GoogleOAuth.ClientID,
		ClientSecret: cfg.GoogleOAuth.ClientSecret,
		RedirectURL:  cfg.GoogleOAuth.RedirectURL,
	}

	authService := services.NewAuthService(logger, jwtHelper, googleConfig, authRepo)
	logger.Info("Init services successfully")

	authHandler := handlers.NewAuthHandler(authService)
	logger.Info("Init handlers successfully")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/google", authHandler.GoogleAuth)
			r.Get("/google/url", authHandler.GoogleAuthURL)
		})
	})

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	logger.Info("starting auth service on %s", cfg.HTTPServer.Address)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server error: ", err)
	}
}
