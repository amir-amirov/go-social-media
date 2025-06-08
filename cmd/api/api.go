package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func newApplication(addr string) *application {
	return &application{
		config: config{
			addr: addr,
		},
	}
}

func (app *application) mount() *chi.Mux {
	// mux := http.NewServeMux()

	// mux.HandleFunc("GET /v1/health", app.healthCheckHandler)

	// return mux

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30, // time it takes for server to send response back to client
		ReadTimeout:  time.Second * 10, // time it takes for server from receiving the first packet to the last packet after tcp connection is formed
		IdleTimeout:  time.Minute,      // time it takes for tcp connection remains alive after sending response back to the client
	}

	log.Printf("Launching server on port%v", app.config.addr)

	return srv.ListenAndServe()
}
