package main

import (
	"context"
	"cursed_backend/internal/config"
	"cursed_backend/internal/db"
	"cursed_backend/internal/logger"
	"cursed_backend/internal/metrics"
	"cursed_backend/internal/router"
	"cursed_backend/internal/security"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  .env not found, using defaults")
	}
	log.Println("✅ .env loaded")

	// Parse and validate config
	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("❌ Config parse failed:", err)
	}
	if cfg.Env == "prod" && cfg.JWTSecret == "" {
		log.Fatal("❌ JWT_SECRET required in prod")
	}
	if cfg.Env == "prod" && cfg.CSRFKey == "" {
		log.Fatal("❌ CSRF_KEY required in prod")
	}
	security.InitJWTSecret(cfg.JWTSecret)
	// Init DB
	db.InitDB(cfg)
	if db.GormDB == nil {
		log.Fatal("❌ GormDB is nil after InitDB()")
	}

	// Init logger
	logger.InitLogger(cfg.LogLevel)
	logger.Log.WithField("config", cfg.Env).Info("App starting")

	// Init metrics
	metrics.InitMetrics()

	// Setup router
	r := router.SetupRouter(cfg)
	logger.Log.Info("✅ Router setup complete")

	// Server setup
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// HTTPS in prod
	if cfg.Env == "prod" {
		go func() {
			if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Log.WithError(err).Fatal("HTTPS server error")
			}
		}()
	} else {
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Log.WithError(err).Fatal("HTTP server error")
			}
		}()
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.WithError(err).Fatal("Server forced to shutdown")
	}

	// Close DB
	if sqlDB, err := db.GormDB.DB(); err == nil {
		err := sqlDB.Close()
		if err != nil {
			logger.Log.WithError(err).Error("Failed to close database connection")
			return
		}
	}
	logger.Log.Info("Server exited cleanly")
}
