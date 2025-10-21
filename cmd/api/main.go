package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/RuLap/meetly-api/meetly/internal/app/auth"
	"github.com/RuLap/meetly-api/meetly/internal/app/event"
	mail_services "github.com/RuLap/meetly-api/meetly/internal/app/mail/services"
	"github.com/RuLap/meetly-api/meetly/internal/app/user"
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

	authModule := auth.NewModule(logger, storage.Database(), jwtHelper, &cfg.GoogleOAuth, redisClient, rabbitmqClient)
	userModule := user.NewModule(logger, storage.Database())

	userProvider := user.NewUserProvider(userModule.Service)

	eventModule := event.NewModule(logger, storage.Database(), userProvider)
	logger.Info("Init modules successfully")

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
	logger.Info("Init mail service successfully")

	router := chi.NewRouter()

	router.Use(chi_middleware.RequestID)
	router.Use(chi_middleware.RealIP)
	router.Use(chi_middleware.Logger)
	router.Use(chi_middleware.Recoverer)
	router.Use(chi_middleware.Timeout(60 * time.Second))

	router.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authModule.Handler.Register)
			r.Post("/login", authModule.Handler.Login)
			r.Post("/google", authModule.Handler.GoogleAuth)
			r.Get("/google/url", authModule.Handler.GoogleAuthURL)
			r.Post("/refresh", authModule.Handler.RefreshTokens)

			r.Route("/email", func(r chi.Router) {
				r.Use(middleware.AuthMiddleware(jwtHelper))

				r.Post("/send-confirmation", authModule.Handler.SendConfirmationLink)
				r.Post("/confirm", authModule.Handler.ConfirmEmail)
			})

			r.With(middleware.AuthMiddleware(jwtHelper)).Post("/logout", authModule.Handler.Logout)
		})
		r.Route("/categories", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtHelper))

			r.Get("/{id}", eventModule.Handler.GetCategoryByID)
			r.Get("/", eventModule.Handler.GetAllCategories)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtHelper))

			r.Get("/{id}", userModule.Handler.GetUserByID)
			r.Put("/{id}", userModule.Handler.UpdateUser)
		})

		r.Route("/events", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtHelper))

			r.Post("/{id}/participants", eventModule.Handler.AddParticipant)

			r.Get("/{id}", eventModule.Handler.GetEventWithDetails)
			r.Get("/", eventModule.Handler.GetShortEvents)
			r.Post("/", eventModule.Handler.CreateEvent)

			r.Get("/categories", eventModule.Handler.GetAllCategories)
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
		logger.Error("server error: ", "error", err)
	}
}

func initRedis(logger *slog.Logger, cfg *config.RedisConfig) *redis.Client {
	logger.Info("starting redis", "address", cfg.Address)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
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
		}
	} else {
		logger.Warn("RabbitMQ URL not configured - email notifications will be disabled")
	}

	return rabbitmqClient
}
