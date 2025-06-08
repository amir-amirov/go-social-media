// It's going to be just a point to check
// if everything is all right with our API and other systems that we consume
package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
