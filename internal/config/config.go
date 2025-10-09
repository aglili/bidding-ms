package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbName    string
	DbUser    string
	DbHost    string
	DbPass    string
	DbPort    string
	AppPort   string
	AppEnv    string
	SecretKey string
	RedisURL  string
	PaystackSecretKey string
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		DbName:    getEnvOrDefault("DB_NAME", ""),
		DbHost:    getEnvOrDefault("DB_HOST", ""),
		DbUser:    getEnvOrDefault("DB_USER", ""),
		DbPass:    getEnvOrDefault("DB_PASSWORD", ""),
		DbPort:    getEnvOrDefault("DB_PORT", ""),
		AppPort:   getEnvOrDefault("APP_PORT", ":6000"),
		AppEnv:    getEnvOrDefault("APP_ENV", "production"),
		SecretKey: getEnvOrDefault("SECRET_KEY", "default_key_trial"),
		RedisURL:  getEnvOrDefault("REDIS_URL", ""),
		PaystackSecretKey: getEnvOrDefault("PAYSTACK_SECRET_KEY",""),
	}
}
