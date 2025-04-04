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
}

func MustGetenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		ConnectionString: MustGetenv("CONNECTION_STRING"),
		TenantID:         MustGetenv("TENANT_ID"),
		ClientID:         MustGetenv("CLIENT_ID"),
	}

	cfg.JWKSURL = "https://login.microsoftonline.com/" + cfg.TenantID + "/discovery/v2.0/keys"

	return cfg
}
