package main

import (
	"context"
	"log"
	"os"

	"iuno-api/db"
	"iuno-api/services/notifications"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env locally if present. Render supplies environment variables
	// directly, so failure to find a local .env file is harmless.
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not configured")
	}

	db.Init(databaseURL)

	if db.Pool == nil {
		log.Fatal("database pool is not initialized")
	}
	defer db.Pool.Close()

	if err := notifications.Run(context.Background()); err != nil {
		log.Fatal("notification job failed: ", err)
	}
}
