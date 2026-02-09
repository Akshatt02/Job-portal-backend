package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(databaseURL string) {
	var err error
	Pool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	log.Println("✅ Connected to PostgreSQL")
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
