package db

import (
	"cursed_backend/models" // Замени на твой модуль, e.g., "zoo-store-backend/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	GormDB *gorm.DB
)

func InitDB() {
	// Загружаем .env файл (если он есть)
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env файл не найден, продолжаем с системными переменными...")
	}

	// Получаем переменные окружения (с дефолтами для локальной разработки)
	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSslmode := os.Getenv("SSL_MODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUsername, dbPassword, dbName, dbPort, dbSslmode)

	// Подключаемся к базе через GORM
	GormDB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("❌ Ошибка подключения к базе данных:", err)
	}

	// Проверяем подключение (ping)
	sqlDB, err := GormDB.DB()
	if err != nil {
		log.Fatal("❌ Ошибка получения sql.DB из GORM:", err)
	}
	if err = sqlDB.Ping(); err != nil {
		log.Fatal("❌ Ошибка пинга базы данных:", err)
	}

	// Автомиграция таблиц (создаст/обновит таблицы, индексы)
	err = GormDB.AutoMigrate(
		&models.User{},    // Таблица users
		&models.Pet{},     // Таблица pets
		&models.Product{}, // Таблица products
	)
	if err != nil {
		log.Fatal("❌ Ошибка миграции таблиц:", err)
	}

	fmt.Println("✅ Таблицы успешно созданы/обновлены и подключение установлено")

}
