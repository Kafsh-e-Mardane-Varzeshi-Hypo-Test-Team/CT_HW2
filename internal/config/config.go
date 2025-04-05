package config

import (
    "fmt"
    "os"
    "strconv"
)

type Config struct {
    Database DatabaseConfig
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DatabaseName string
    SSLMode  string
}

// ConnectionString returns a formatted PostgreSQL connection string
func (c DatabaseConfig) ConnectionString() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=%s",
        c.User, c.Password, c.Host, c.Port, c.DatabaseName, c.SSLMode,
    )
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    portStr := os.Getenv("DB_PORT")
    if portStr == "" {
        return nil, fmt.Errorf("DB_PORT environment variable is not set")
    }
    
    dbPort, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("DB_PORT must be a valid number: %w", err)
    }

    host := os.Getenv("DB_HOST")
    if host == "" {
        return nil, fmt.Errorf("DB_HOST environment variable is not set")
    }

    user := os.Getenv("DB_USER")
    if user == "" {
        return nil, fmt.Errorf("DB_USER environment variable is not set")
    }

    password := os.Getenv("DB_PASSWORD")
    if password == "" {
        return nil, fmt.Errorf("DB_PASSWORD environment variable is not set")
    }

    databaseName := os.Getenv("DB_NAME")
    if databaseName == "" {
        return nil, fmt.Errorf("DB_NAME environment variable is not set")
    }

    sslMode := os.Getenv("DB_SSLMODE")
    if sslMode == "" {
        return nil, fmt.Errorf("DB_SSLMODE environment variable is not set")
    }
    
    return &Config{
        Database: DatabaseConfig{
            Host:     host,
            Port:     dbPort,
            User:     user,
            Password: password,
            DatabaseName: databaseName,
            SSLMode:  sslMode,
        },
    }, nil
}