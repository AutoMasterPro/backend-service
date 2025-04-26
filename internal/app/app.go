package app

import (
	"backend-service/internal/config"
	"backend-service/internal/handlers"
	"backend-service/internal/services"
	"backend-service/internal/storages"
	"backend-service/pkg/database"
	"backend-service/pkg/jwt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func Run() {
	// logger
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"}
	logger := zerolog.New(output).With().Caller().Timestamp().Logger()
	dir, err := os.Getwd()
	if err != nil {
		logger.Error().Msg("Cannot get working directory")
	}
	logger.Info().Msg("Current working directory: " + dir)
	// init env
	err = godotenv.Load()
	if err != nil {
		logger.Error().Msgf("Error loading .env file: %v", err)
	}
	// cfg
	cfg := config.GetConfig()
	logger.Info().Msg("Config: OK")
	// postgres
	pg, err := database.NewPostgresDB(cfg.Postgres.DBHost, cfg.Postgres.DBPort, cfg.Postgres.DBUser, cfg.Postgres.DBName, cfg.Postgres.DBPass, cfg.Postgres.DBSSLMode)
	if err != nil {
		logger.Error().Msgf("Error connecting to PostgreSQL: %v", err)
	}
	logger.Info().Msg("Postgres: OK")
	// storage
	storage := storages.NewStorage(storages.StorageDeps{
		PostgresDB: pg,
		Log:        logger,
	})
	// services
	service := services.NewService(services.ServiceDeps{
		Log:     logger,
		Storage: storage,
	})
	// jwt service
	jwtService := jwt.New(jwt.Config{
		SecretKey:       cfg.AppSecretKey,
		AccessTokenTTL:  time.Hour * 24 * 7,
		RefreshTokenTTL: 0,
	}, nil)
	// handlers
	handler := handlers.NewHandler(logger, service, jwtService)
	// run
	handler.InitRoutes(cfg.AppPort)
}
