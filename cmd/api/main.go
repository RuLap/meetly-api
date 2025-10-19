package main

import (
	"context"
	"encoding/json"
	auth_handlers "github.com/RuLap/meetly-api/meetly/internal/app/auth/handlers"
	auth_repos "github.com/RuLap/meetly-api/meetly/internal/app/auth/repository"
	auth_services "github.com/RuLap/meetly-api/meetly/internal/app/auth/services"
	category_handlers "github.com/RuLap/meetly-api/meetly/internal/app/category/handlers"
	category_repos "github.com/RuLap/meetly-api/meetly/internal/app/category/repository"
	category_services "github.com/RuLap/meetly-api/meetly/internal/app/category/services"
	mail_services "github.com/RuLap/meetly-api/meetly/internal/app/mail/services"
	user_handlers "github.com/RuLap/meetly-api/meetly/internal/app/user/handlers"
	user_repos "github.com/RuLap/meetly-api/meetly/internal/app/user/repository"
	user_services "github.com/RuLap/meetly-api/meetly/internal/app/user/services"
	"github.com/RuLap/meetly-api/meetly/internal/pkg/config"
	"github.com/RuLap/meetly-api/meetly/internal/pkg/jwt_helper"
	"github.com/RuLap/meetly-api/meetly/internal/pkg/logger"
	"github.com/RuLap/meetly-api/meetly/internal/pkg/middleware"
	"github.com/RuLap/meetly-api/meetly/internal/pkg/rabbitmq"
	postgres "github.com/RuLap/meetly-api/meetly/internal/pkg/storage"
	validation "github.com/RuLap/meetly-api/meetly/internal/pkg/validator"
	"github.com/darahayes/go-boom"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.New(logger.Config{
		Level:   cfg.Env,
		LokiURL: cfg.Log.LokiURL,
		Labels:  cfg.Log.LokiLabels,
	})

	validation.Init()

	redisClient := initRedis(logger, &cfg.Redis)
	logger.Info("init redis successfully")

	rabbitmqClient := initRabbitMQ(logger, &cfg.RabbitMQ)
	logger.Info("init rabbitmq successfully")

	storage, err := postgres.InitDB(cfg.PostgresConnString)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		return
	}

	jwtHelper, err := jwt_helper.NewJwtHelper(cfg.JWT.Secret)
	if err != nil {
		logger.Error("failed to create JWT helper", "error", err)
		return
	}

	authRepo := auth_repos.NewAuthRepository(storage.Database())
	categoryRepo := category_repos.NewCategoryRepository(storage.Database())
	userRepo := user_repos.NewUserRepository(storage.Database())
	logger.Info("Init repos successfully")

	googleConfig := &auth_services.GoogleOAuthConfig{
		ClientID:     cfg.GoogleOAuth.ClientID,
		ClientSecret: cfg.GoogleOAuth.ClientSecret,
		RedirectURL:  cfg.GoogleOAuth.RedirectURL,
	}

	authService := auth_services.NewAuthService(logger, jwtHelper, googleConfig, redisClient, rabbitmqClient, authRepo)
	categoryService := category_services.NewCategoryService(logger, categoryRepo)
	userService := user_services.NewUserService(logger, userRepo)
	var mailService *mail_services.MailService
	if rabbitmqClient != nil {
		mailService = mail_services.NewMailService(
			logger,
			rabbitmqClient,
			&cfg.SMTP,
		)

		go func() {
			logger.Info("starting mail service consumer")
			if err := mailService.StartConsumer(context.Background()); err != nil {
				logger.Error("mail service consumer failed", "error", err)
			}
		}()
	} else {
		logger.Warn("mail service not started - RabbitMQ not available")
	}
	logger.Info("Init services successfully")

	authHandler := auth_handlers.NewAuthHandler(authService)
	categoryHandler := category_handlers.NewCategoryHandler(categoryService)
	userHandler := user_handlers.NewUserHandler(userService)
	logger.Info("Init handlers successfully")

	router := chi.NewRouter()

	router.Use(chi_middleware.RequestID)
	router.Use(chi_middleware.RealIP)
	router.Use(chi_middleware.Logger)
	router.Use(chi_middleware.Recoverer)
	router.Use(chi_middleware.Timeout(60 * time.Second))

	router.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/google", authHandler.GoogleAuth)
			r.Get("/google/url", authHandler.GoogleAuthURL)
			r.Post("/refresh", authHandler.RefreshTokens)

			r.Route("/email", func(r chi.Router) {
				r.Use(middleware.AuthMiddleware(jwtHelper))

				r.Post("/send-confirmation", authHandler.SendConfirmationLink)
				r.Post("/confirm", authHandler.ConfirmEmail)
			})

			r.With(middleware.AuthMiddleware(jwtHelper)).Post("/logout", authHandler.Logout)
		})
		r.Route("/categories", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtHelper))

			r.Get("/{id}", categoryHandler.GetCategoryByID)
			r.Get("", categoryHandler.GetAllCategories)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtHelper))

			r.Get("/{id}", userHandler.GetUserByID)
			r.Put("/{id}", userHandler.UpdateUser)
		})
	})

	router.Get("/health/rabbitmq", func(w http.ResponseWriter, r *http.Request) {
		if rabbitmqClient == nil {
			boom.ServerUnavailable(w, "RabbitMQ not configured")
			return
		}

		if err := rabbitmqClient.HealthCheck(); err != nil {
			boom.ServerUnavailable(w, "RabbitMQ connection failed")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "healthy",
			"service": "rabbitmq",
		})
	})

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	logger.Info("starting", "address", cfg.HTTPServer.Address)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server error: ", err)
	}
}

func initRedis(logger *slog.Logger, cfg *config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Error("failed to connect to redis", "error", err)
	}

	logger.Info("redis connected successfully")
	return rdb
}

func initRabbitMQ(logger *slog.Logger, cfg *config.RabbitMQConfig) *rabbitmq.Client {
	var rabbitmqClient *rabbitmq.Client
	var err error

	if cfg.URL != "" {
		rabbitmqClient, err = rabbitmq.NewClient(cfg.URL, logger)
		if err != nil {
			logger.Error("failed to connect to RabbitMQ",
				"error", err,
				"url", cfg.URL,
			)
		} else {
			logger.Info("successfully connected to RabbitMQ")

			defer func() {
				if rabbitmqClient != nil {
					if err := rabbitmqClient.Close(); err != nil {
						logger.Error("failed to close RabbitMQ connection", "error", err)
					} else {
						logger.Info("RabbitMQ connection closed")
					}
				}
			}()
		}
	} else {
		logger.Warn("RabbitMQ URL not configured - email notifications will be disabled")
	}

	return rabbitmqClient
}
