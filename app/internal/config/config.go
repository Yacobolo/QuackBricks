package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ConnectionString string
	TenantID         string
	ClientID         string
	JWKSURL          string
	Endpoint         string
	Scopes           []string
}

func mustGetEnv(key string) string {
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
		ConnectionString: mustGetEnv("CONNECTION_STRING"),
		TenantID:         mustGetEnv("TENANT_ID"),
		ClientID:         mustGetEnv("CLIENT_ID"),
		Endpoint:         getEnv("ENDPOINT", "http://localhost:8080/api"),
	}

	cfg.JWKSURL = "https://login.microsoftonline.com/" + cfg.TenantID + "/discovery/v2.0/keys"

	cfg.Scopes = []string{cfg.ClientID + "/.default"}

	return cfg
}
