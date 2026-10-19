package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Environment   string
	DatabaseUrl   string
	AdminPassword string
	SessionSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:          toGoPort(getOrDefault("PORT", ":8080")),
		Environment:   getOrDefault("ENVIRONMENT", "production"),
		DatabaseUrl:   requrie("DATABASE_URL"),
		AdminPassword: requrie("ADMIN_PASSWORD"),
		SessionSecret: requrie("SESSION_SECRET"),
	}
}

func getOrDefault(key string, defaltV string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaltV
	}
	return v
}

func toGoPort(port string) string {
	if port == "" {
		return ":8080"
	}
	if port[0] != ':' {
		return ":" + port
	}
	return port
}

func requrie(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required environment variable '%s'", key)
	}
	return v
}
