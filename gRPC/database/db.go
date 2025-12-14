package database

import (
	"Auth/models"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() error {
	// Ensure the directory for the DB exists
	dbPath := filepath.Join("data", "storage.db") // Better path organization
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
		return err
	}

	// Open or create database
	var err error
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}

	// Run migrations for all your models
	err = db.AutoMigrate(
		&models.LogData{},
		&models.Entry{}, // Make sure to include your Entry model
		// Add other models here
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
		return err
	}

	log.Println("Database setup complete with all tables")
	return nil
}

func GetDB() *gorm.DB {
	return db
}

func Close() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting generic database object: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}
