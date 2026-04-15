package main

import (
	"log"
	"os"

	"learning-growth-platform/internal/config"
	"learning-growth-platform/internal/database"
	"learning-growth-platform/internal/http/router"
)

func main() {
	_ = config.LoadDotEnv(".env")
	cfg := config.Load()
	if cfg.MySQLDSN == "" {
		log.Fatal("MYSQL_DSN is required; please set it in environment or .env")
	}
	if cfg.JWTSecret != "" {
		_ = os.Setenv("AUTH_TOKEN_SECRET", cfg.JWTSecret)
	}

	db, err := database.OpenMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate mysql: %v", err)
	}

	r := router.NewRouter(db)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
