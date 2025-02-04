package inits

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/XEFF09/fulldev-be/cogs/entities"
)

func DatabaseConnection() *gorm.DB {

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	DB_HOST := os.Getenv("DATABASE_HOST")	
	DB_USER := os.Getenv("DATABASE_USER")
	DB_PASSWORD := os.Getenv("DATABASE_PASSWORD")
	DB_NAME := os.Getenv("DATABASE_NAME")

	DB_PORT, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		panic("Error could not convert env data type")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to database!")
	}
	fmt.Println("Database connected!")

	err = db.AutoMigrate(entities.User{})
	if err != nil {
		panic("Failed to migrate database!")
	}
	fmt.Println("Database migrated!")

	return db
}