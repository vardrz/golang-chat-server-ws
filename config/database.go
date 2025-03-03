package config

import (
	"chat-server/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	
	fmt.Println(dsn)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        fmt.Println("failed to connect database")
        log.Fatal(err)
    }

    return db
}

func InitMigration(db *gorm.DB) {
    // Print current tables
    tables, err := db.Migrator().GetTables()
    if err != nil {
        log.Fatal("Failed to get tables:", err)
    }
    log.Println("Existing tables:", tables)	

    // Auto migrate models with detailed error logging
    err = db.AutoMigrate(&models.User{}, &models.Message{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    // Verify if tables exist
    if !db.Migrator().HasTable(&models.Message{}) {
        log.Fatal("Messages table was not created")
    } else {
        log.Println("Messages table created successfully")
    }
}