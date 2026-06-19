package database

import (
	"log"
	"pw_film/internal/config"
	"pw_film/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB connects to the PostgreSQL database using GORM and runs auto-migrations
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database.")

	// Automatically run migrations for registered models
	if err := db.AutoMigrate(&models.Film{}); err != nil {
		return nil, err
	}
	log.Println("Database migration completed.")

	return db, nil
}
