package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	API      APIConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port     string
	Mode     string
	LogLevel string
}

type APIConfig struct {
	CoinGeckoURL    string
	UpdateInterval  int
	RequestTimeout  int
}

func Load() *Config {
	// Загружаем .env файл если он существует
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "exchange_user"),
			Password: getEnv("DB_PASSWORD", "exchange_password"),
			Name:     getEnv("DB_NAME", "crypto_exchange"),
		},
		Server: ServerConfig{
			Port:     getEnv("PORT", "8080"),
			Mode:     getEnv("GIN_MODE", "debug"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		API: APIConfig{
			CoinGeckoURL:    getEnv("COINGECKO_URL", "https://api.coingecko.com/api/v3"),
			UpdateInterval:  getEnvAsInt("API_UPDATE_INTERVAL", 60),
			RequestTimeout:  getEnvAsInt("REQUEST_TIMEOUT", 10),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}