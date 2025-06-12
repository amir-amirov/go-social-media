// It's going to be just a point to check
// if everything is all right with our API and other systems that we consume
package main

import (
	"log"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		// error
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		log.Printf("[ERROR] %v", err.Error())
	}
}
