package config

import "os"

type Config struct {
	Port      string
	MySQLDSN  string
	JWTSecret string
}

func Load() Config {
	return Config{
		Port: getEnvOrDefault("APP_PORT", "8080"),
		MySQLDSN: getEnvOrDefault(
			"MYSQL_DSN",
			"root:010511@tcp(127.0.0.1:3306)/learning_growth?charset=utf8mb4&parseTime=True&loc=Local",
		),
		JWTSecret: getEnvOrDefault("JWT_SECRET", "local-dev-secret"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
