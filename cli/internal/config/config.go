package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TenantID string
	ClientID string
	Endpoint string
	Scopes   []string
}

func MustGetenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		TenantID: MustGetenv("TENANT_ID"),
		ClientID: MustGetenv("CLIENT_ID"),
		Endpoint: getEnv("ENDPOINT", "http://localhost:8080/api"),
	}
	cfg.Scopes = []string{cfg.ClientID + "/.default"}

	return cfg
}
