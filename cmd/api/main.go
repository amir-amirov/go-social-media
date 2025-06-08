package main

import "log"

func main() {
	app := newApplication(":8080")

	if err := app.run(); err != nil {
		log.Fatalf("[ERROR] Unable to run application..")
	}
}
