package main

import (
	"log"

	"github.com/amir-amirov/go-social-media/internal/env"
	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	store := store.NewPostgresStorage(nil)

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		},
	}

	app := newApplication(cfg, store)

	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatalf("[ERROR] Unable to launch server..")
	}
}
