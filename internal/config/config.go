package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	APIKey          string
	BaseURL         string
	LogLevel        string
	CacheExpiration int
	UpdateInterval  int
	MaxHistoryDays  int
}

func Load() (*Config, error) {
	godotenv.Load()

	config := &Config{
		Port:            getEnv("PORT", "8080"),
		APIKey:          getEnv("EXCHANGE_API_KEY", "73fda159574e487335c01f410f990937"),
		BaseURL:         getEnv("EXCHANGE_BASE_URL", "https://api.exchangerate.host"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		CacheExpiration: getEnvAsInt("CACHE_EXPIRATION", 3600),
		UpdateInterval:  getEnvAsInt("UPDATE_INTERVAL", 14400),
		MaxHistoryDays:  getEnvAsInt("MAX_HISTORY_DAYS", 90),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
