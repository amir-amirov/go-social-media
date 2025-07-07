package main

import (
	"log"

	"github.com/amir-amirov/go-social-media/internal/db"
	"github.com/amir-amirov/go-social-media/internal/env"
	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := env.GetString("SEED_DB_ADDR", "")
	conn, err := db.New(addr, 3, 3)
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewPostgresStorage(conn)

	db.Seed(store)
}
