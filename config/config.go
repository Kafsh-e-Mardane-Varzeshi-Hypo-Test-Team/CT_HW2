package config

import (
	"log"

	"github.com/joho/godotenv"
)

// loads environment variables from a .env file
// TODO: this function should be replaced with docker
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}
}
