package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file when present.
// In container deployments where variables are injected directly,
// the absence of a .env file is not treated as a fatal error.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}
}
