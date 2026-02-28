package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	MongoURI       string
	DBName         string
	JWTSecret      string
	JWTExpiryHours int
}

// LoadEnv loads config from environment. Call Validate() after load.
func LoadEnv() *Config {
	hours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if hours <= 0 {
		hours = 24
	}
	return &Config{
		Port:           getEnv("PORT", "8080"),
		MongoURI:       getEnv("MONGO_URI", ""),
		DBName:         getEnv("DB_NAME", "planelite"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiryHours: hours,
	}
}

// Validate fails fast on missing required config. Call after LoadEnv().
func (c *Config) Validate() error {
	if c.MongoURI == "" {
		return fmt.Errorf("config: MONGO_URI is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("config: JWT_SECRET is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
