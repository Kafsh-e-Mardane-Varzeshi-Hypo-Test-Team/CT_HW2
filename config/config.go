package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
)

type ServerConfig struct {
	Host string
	Port int
}

type JWTConfig struct {
	SecretKey string
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	SSLMode      string
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// ConnectionString returns a formatted PostgreSQL connection string
func (c DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DatabaseName, c.SSLMode,
	)
}

// LoadConfig loads configuration from environment variables
func LoadConfig() error {
	var err error
	// TODO: docker-compose up
	loadEnv()

	err = LoadServerConfig()
	if err != nil {
		return fmt.Errorf("failed to load server config: %w", err)
	}

	err = LoadJWTConfig()
	if err != nil {
		return fmt.Errorf("failed to load JWT config: %w", err)
	}

	err = LoadDatabaseConfig()
	if err != nil {
		return fmt.Errorf("failed to load database config: %w", err)
	}

	return nil
}

func LoadServerConfig() error {
	host, err := getEnv("SERVER_HOST")
	if err != nil {
		return err
	}

	portStr, err := getEnv("SERVER_PORT")
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("SERVER_PORT must be a valid number: %w", err)
	}

	Server = ServerConfig{
		Host: host,
		Port: port,
	}
	return nil
}

func LoadJWTConfig() error {
	js, err := getEnv("JWT_SECRET")
	if err != nil {
		return err
	}

	jwtSecret, err := base64.StdEncoding.DecodeString(js)
	if err != nil {
		return fmt.Errorf("invalid JWT_SECRET: %w", err)
	}

	JWT = JWTConfig{
		SecretKey: string(jwtSecret),
	}
	return nil
}

func LoadDatabaseConfig() error {
	portStr, err := getEnv("DB_PORT")
	if err != nil {
		return err
	}

	dbPort, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("DB_PORT must be a valid number: %w", err)
	}

	host, err := getEnv("DB_HOST")
	if err != nil {
		return err
	}

	user, err := getEnv("DB_USER")
	if err != nil {
		return err
	}

	password, err := getEnv("DB_PASSWORD")
	if err != nil {
		return err
	}

	databaseName, err := getEnv("DB_NAME")
	if err != nil {
		return err
	}

	sslMode, err := getEnv("DB_SSLMODE")
	if err != nil {
		return err
	}

	Database = DatabaseConfig{
		Host:         host,
		Port:         dbPort,
		User:         user,
		Password:     password,
		DatabaseName: databaseName,
		SSLMode:      sslMode,
	}
	return nil
}

// Gets an environment variable or returns an error
func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is required", key)
	}
	return value, nil
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}
}
