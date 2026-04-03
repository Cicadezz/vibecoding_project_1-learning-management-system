package config

import "os"

type Config struct {
	Port      string
	MySQLDSN  string
	JWTSecret string
}

func Load() Config {
	return Config{
		Port:      getEnvOrDefault("APP_PORT", "8080"),
		MySQLDSN:  os.Getenv("MYSQL_DSN"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
