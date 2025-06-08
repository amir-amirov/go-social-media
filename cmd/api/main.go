package main

import (
	"log"

	"github.com/amir-amirov/go-social-media/internal/env"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := newApplication(cfg)
	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatalf("[ERROR] Unable to launch server..")
	}
}
