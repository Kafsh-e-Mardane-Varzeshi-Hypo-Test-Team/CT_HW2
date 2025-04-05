package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey string
}

// ConnectionString returns a formatted PostgreSQL connection string
func (c DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
}

var Database DatabaseConfig
var JWT JWTConfig

// LoadConfig loads configuration from environment variables
func LoadConfig() error {
	// TODO: docker-compose up
	// loadEnv()

	// database configuration
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return fmt.Errorf("invalid DB_PORT: %w", err)
	}

	Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Database: getEnv("DB_NAME", "hw2"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// JWT configuration
	js := getEnv("JWT_SECRET", "")
	jwtSecret, err := base64.StdEncoding.DecodeString(js)

	if err != nil {
		return fmt.Errorf("invalid JWT_SECRET: %w", err)
	}

	JWT = JWTConfig{
		SecretKey: string(jwtSecret),
	}

	return nil
}

// Gets an environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}
}
