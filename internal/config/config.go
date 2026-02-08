package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Load() *Config {
	_ = godotenv.Load() // load .env, ignore error nếu chạy production

	return &Config{
		App: AppConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     getEnv("DB_NAME", "pingpong"),
		},
	}
}
