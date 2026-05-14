package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init(connStr string) {
	var err error

	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	log.Println("Connected to PostgreSQL")
}