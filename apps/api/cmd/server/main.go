package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/foundryhq/foundryhq/apps/api/internal/config"
	"github.com/foundryhq/foundryhq/apps/api/internal/handlers"
	"github.com/foundryhq/foundryhq/apps/api/pkg/database"
	"github.com/foundryhq/foundryhq/apps/api/pkg/logger"
)

func main() {
	// Config must load before the structured logger exists, so failures here
	// use the standard log package rather than zap.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	zapLogger, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatalf("creating logger: %v", err)
	}
	// zap buffers writes for performance; Sync flushes them so logs aren't
	// lost if the process exits abruptly.
	defer zapLogger.Sync()

	db, err := database.Connect(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		zapLogger.Fatal("connecting to database", zap.Error(err))
	}

	// ReleaseMode silences gin's debug-level request logging/warnings, which
	// are noisy and unnecessary once structured logging is in place.
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	// gin.New (not gin.Default) skips gin's own logger middleware since zap
	// already covers logging; Recovery is added back explicitly so a panic
	// in one handler returns a 500 instead of crashing the whole server.
	router := gin.New()
	router.Use(gin.Recovery())

	healthHandler := handlers.NewHealthHandler(db)
	router.GET("/health", healthHandler.Check)

	zapLogger.Info("starting server", zap.String("port", cfg.Port), zap.String("env", cfg.Env))
	// Run blocks, serving until the process is killed or a listener error
	// occurs; the leading ":" binds to all network interfaces.
	if err := router.Run(":" + cfg.Port); err != nil {
		zapLogger.Fatal("server failed", zap.Error(err))
	}
}
