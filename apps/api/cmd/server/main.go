package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/foundryhq/foundryhq/apps/api/internal/config"
	"github.com/foundryhq/foundryhq/apps/api/internal/handlers"
	"github.com/foundryhq/foundryhq/apps/api/internal/middleware"
	"github.com/foundryhq/foundryhq/apps/api/internal/repositories/postgres"
	"github.com/foundryhq/foundryhq/apps/api/internal/usecases"
	"github.com/foundryhq/foundryhq/apps/api/pkg/database"
	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
	"github.com/foundryhq/foundryhq/apps/api/pkg/logger"
)

// shutdownTimeout bounds how long graceful shutdown waits for in-flight
// requests to finish before forcing the listener closed.
const shutdownTimeout = 10 * time.Second

func main() {
	// Config must load before the structured logger exists, so failures here
	// use the standard log package rather than zap.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}
	fmt.Println("✓ Configuration loaded")

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
	fmt.Println("✓ PostgreSQL connected")

	// ReleaseMode silences gin's debug-level request logging/warnings, which
	// are noisy and unnecessary once structured logging is in place.
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	// gin.New (not gin.Default) skips gin's own logger/recovery middleware —
	// both are replaced below with zap-backed equivalents. Recovery goes
	// first so it wraps (and can catch panics from) every middleware after
	// it; RequestID goes before Logger so request IDs are available to log.
	router := gin.New()
	router.Use(
		middleware.Recovery(zapLogger),
		middleware.RequestID(),
		middleware.Logger(zapLogger),
		middleware.CORS(cfg.CORSAllowedOrigins),
	)

	healthHandler := handlers.NewHealthHandler(db)
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/version", healthHandler.Version)

	jwtManager := jwt.NewManager(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)
	authUsecase := usecases.NewAuthUsecase(
		postgres.NewUserRepository(db),
		postgres.NewRefreshTokenRepository(db),
		jwtManager,
	)
	// secureCookies (the refresh-token cookie's Secure flag) must be false
	// for local HTTP dev — browsers refuse Secure cookies over plain HTTP —
	// and true in production, which is expected to run behind HTTPS.
	authHandler := handlers.NewAuthHandler(authUsecase, cfg.Env == "production")
	// loginRateLimiter throttles brute-force credential guessing per client
	// IP — burst requests allowed immediately, refilling gradually over the
	// configured window rather than resetting all at once.
	loginRateLimiter := middleware.NewRateLimiter(
		rate.Every(cfg.LoginRateLimitWindow/time.Duration(cfg.LoginRateLimitBurst)),
		cfg.LoginRateLimitBurst,
	)
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", loginRateLimiter.Limit(), authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.Logout)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// ListenAndServe blocks, so it runs in its own goroutine; the main
	// goroutine instead waits below for an OS signal to begin shutdown.
	go func() {
		fmt.Printf("✓ Server running on :%s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zapLogger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("shutting down server")

	// Bounding shutdown with a timeout stops it from hanging forever on a
	// stuck connection; requests still in flight get this long to finish.
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("server forced to shut down", zap.Error(err))
	}

	if sqlDB, err := db.DB(); err != nil {
		zapLogger.Error("unwrapping database handle", zap.Error(err))
	} else if err := sqlDB.Close(); err != nil {
		zapLogger.Error("closing database connection", zap.Error(err))
	}

	zapLogger.Info("server exited")
}
