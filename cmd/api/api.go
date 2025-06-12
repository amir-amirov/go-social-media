package main

import (
	"log"
	"net/http"
	"time"

	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
}

func newApplication(cfg config, store store.Storage) *application {
	return &application{
		config: cfg,
		store:  store,
	}
}

func (app *application) mount() http.Handler {
	// mux := http.NewServeMux()

	// mux.HandleFunc("GET /v1/health", app.healthCheckHandler)

	// return mux

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {

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
