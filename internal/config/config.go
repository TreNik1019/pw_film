package config

import (
	"fmt"
	"os"
)

// Config holds all configurations for the application
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSslMode  string
	ServerPort string
}

// LoadConfig initializes the configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "film"),
		DBSslMode:  getEnv("DB_SSLMODE", "disable"),
		ServerPort: getEnv("PORT", "8080"),
	}
}

// GetDSN formats the connection string for GORM PostgreSQL driver
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSslMode)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
