package db

import (
	"cursed_backend/internal/config"
	"cursed_backend/internal/logger"
	"cursed_backend/internal/metrics"
	"cursed_backend/internal/models"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

type dbQueryHook struct {
	start time.Time
}

func (h *dbQueryHook) Name() string { return "dbQueryHook" }

func (h *dbQueryHook) Initialize(db *gorm.DB) error {
	logger.Log.Info("Initializing database query hooks")

	err := db.Callback().Create().Before("gorm:before_create").Register("query_start", h.beforeQuery)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to register Create.Before callback")
		return err
	}
	logger.Log.Debug("Registered Create.Before callback")

	err = db.Callback().Create().After("gorm:after_create").Register("query_metrics", h.afterQuery)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to register Create.After callback")
		return err
	}
	logger.Log.Debug("Registered Create.After callback")

	err = db.Callback().Query().Before("gorm:before_query").Register("query_start", h.beforeQuery)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to register Query.Before callback")
		return err
	}
	logger.Log.Debug("Registered Query.Before callback")

	err = db.Callback().Query().After("gorm:after_query").Register("query_metrics", h.afterQuery)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to register Query.After callback")
		return err
	}
	logger.Log.Debug("Registered Query.After callback")

	err = db.Callback().Update().Before("gorm:before_update").Register("query_start", h.beforeQuery)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to register Update.Before callback")
		return err
	}
	logger.Log.Debug("Registered Update.Before callback")

	err = db.Callback().Update().After("gorm:after_update").Register("query_metrics", h.afterQuery)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to register Update.After callback")
		return err
	}
	logger.Log.Debug("Registered Update.After callback")

	err = db.Callback().Delete().Before("gorm:before_delete").Register("query_start", h.beforeQuery)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to register Delete.Before callback")
		return err
	}
	logger.Log.Debug("Registered Delete.Before callback")

	err = db.Callback().Delete().After("gorm:after_delete").Register("query_metrics", h.afterQuery)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to register Delete.After callback")
		return err
	}
	logger.Log.Debug("Registered Delete.After callback")

	logger.Log.Info("Database query hooks initialized successfully")
	return nil
}

func (h *dbQueryHook) beforeQuery(*gorm.DB) {
	h.start = time.Now()
	logger.Log.Debug("Query started")
}

func (h *dbQueryHook) afterQuery(db *gorm.DB) {
	table := "unknown"
	if db.Statement != nil && db.Statement.Table != "" {
		table = db.Statement.Table
	}
	duration := time.Since(h.start).Seconds()
	metrics.DBQueryDuration.WithLabelValues(table).Observe(duration)
	logger.Log.WithFields(map[string]interface{}{
		"table":    table,
		"duration": duration,
	}).Debug("Query completed")
}

func InitDB(cfg *config.Config) {
	logger.Log.Info("Initializing database connection")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUsername, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.SSLMode)

	if cfg.DBHost == "" || cfg.DBUsername == "" || cfg.DBName == "" {
		logger.Log.Fatal("Missing required DB env vars")
	}

	logger.Log.WithFields(map[string]interface{}{
		"host":   cfg.DBHost,
		"dbname": cfg.DBName,
		"port":   cfg.DBPort,
	}).Info("Connecting to database")

	var err error
	GormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
	})
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to open GORM connection")
	}
	logger.Log.Info("GORM connection established")

	err = GormDB.Use(&dbQueryHook{})
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to register database query hooks")
	}

	sqlDB, err := GormDB.DB()
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to get underlying sql.DB")
	}

	logger.Log.WithFields(map[string]interface{}{
		"max_open_conns": 25,
		"max_idle_conns": 10,
	}).Info("Configuring connection pool")
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	logger.Log.Info("Pinging database")
	if err = sqlDB.Ping(); err != nil {
		logger.Log.WithError(err).Fatal("Failed to ping database")
	}
	logger.Log.Info("Database ping successful")

	logger.Log.Info("Running database migrations")
	if err = GormDB.AutoMigrate(&models.User{}, &models.Pet{}, &models.Product{}); err != nil {
		logger.Log.WithError(err).Fatal("Failed to run migrations")
	}
	logger.Log.Info("Database migrations completed")

	logger.Log.Info("Database connected and migrated successfully")
}
