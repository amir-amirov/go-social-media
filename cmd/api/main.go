package main

import (
	"log"
)

func main() {
	app := newApplication(":3000")
	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatalf("[ERROR] Unable to launch server..")
	}
}
