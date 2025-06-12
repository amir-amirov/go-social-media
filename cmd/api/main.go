package main

import (
	"log"

	"github.com/amir-amirov/go-social-media/internal/db"
	"github.com/amir-amirov/go-social-media/internal/env"
	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost:5431/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns)
	if err != nil {
		log.Fatalf("[ERROR] Unable to connect to database..")
	}

	defer db.Close()

	store := store.NewPostgresStorage(db)

	app := newApplication(cfg, store)

	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatalf("[ERROR] Unable to launch server..")
	}
}
