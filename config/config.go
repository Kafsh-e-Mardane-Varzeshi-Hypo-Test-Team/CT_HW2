package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type JWTConfig struct {
	SecretKey []byte
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	SSLMode      string
}

func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// ConnectionString returns a formatted PostgreSQL connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DatabaseName, c.SSLMode,
	)
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// TODO: docker-compose up
	loadEnv()

	serverConfig, err := LoadServerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load server config: %w", err)
	}

	jwtConfig, err := LoadJWTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load JWT config: %w", err)
	}

	databaseConfig, err := LoadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	return &Config{
		Server:   serverConfig,
		JWT:      jwtConfig,
		Database: databaseConfig,
	}, nil
}

func LoadServerConfig() (ServerConfig, error) {
	host, err := getEnv("SERVER_HOST")
	if err != nil {
		return ServerConfig{}, err
	}

	portStr, err := getEnv("SERVER_PORT")
	if err != nil {
		return ServerConfig{}, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("SERVER_PORT must be a valid number: %w", err)
	}

	return ServerConfig{
		Host: host,
		Port: port,
	}, nil
}

func LoadJWTConfig() (JWTConfig, error) {
	js, err := getEnv("JWT_SECRET")
	if err != nil {
		return JWTConfig{}, err
	}

	jwtSecret, err := base64.StdEncoding.DecodeString(js)
	if err != nil {
		return JWTConfig{}, fmt.Errorf("invalid JWT_SECRET: %w", err)
	}

	return JWTConfig{
		SecretKey: jwtSecret,
	}, nil
}

func LoadDatabaseConfig() (DatabaseConfig, error) {
	portStr, err := getEnv("DB_PORT")
	if err != nil {
		return DatabaseConfig{}, err
	}

	dbPort, err := strconv.Atoi(portStr)
	if err != nil {
		return DatabaseConfig{}, fmt.Errorf("DB_PORT must be a valid number: %w", err)
	}

	host, err := getEnv("DB_HOST")
	if err != nil {
		return DatabaseConfig{}, err
	}

	user, err := getEnv("DB_USER")
	if err != nil {
		return DatabaseConfig{}, err
	}

	password, err := getEnv("DB_PASSWORD")
	if err != nil {
		return DatabaseConfig{}, err
	}

	databaseName, err := getEnv("DB_NAME")
	if err != nil {
		return DatabaseConfig{}, err
	}

	sslMode, err := getEnv("DB_SSLMODE")
	if err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		Host:         host,
		Port:         dbPort,
		User:         user,
		Password:     password,
		DatabaseName: databaseName,
		SSLMode:      sslMode,
	}, nil
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
