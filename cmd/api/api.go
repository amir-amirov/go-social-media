package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amir-amirov/go-social-media/docs"
	"github.com/amir-amirov/go-social-media/internal/mailer"
	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
	mailer mailer.Client
}

type config struct {
	addr   string
	apiURL string
	db     dbConfig
	env    string
	mail   mailConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
}

type mailConfig struct {
	exp       time.Duration
	sendgrid  sendGridConfig
	mailtrap  mailTrapConfig
	fromEmail string
}

type sendGridConfig struct {
	apiKey string
}

type mailTrapConfig struct {
	apiKey string
}

func newApplication(cfg config, store store.Storage, logger *zap.SugaredLogger, mailer mailer.Client) *application {
	log.Println("cfg:", cfg)
	// log.Println("store:", store)
	// log.Println("logger: ", logger)
	log.Println("mailer: ", mailer)
	return &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
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

	docURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)

	r.Route("/v1", func(r chi.Router) {

		// r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)
		r.Get("/health", app.healthCheckHandler)

		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Delete("/", app.checkPostOwnership("moderator", app.deletePostHandler))

				r.Route("/comments", func(r chi.Router) {

					r.Post("/", app.createCommentHandler)

					r.Route("/{commentID}", func(r chi.Router) {
						r.Use(app.commentsContextMiddleware)

						r.Get("/", app.getCommentsHandler)
						// TODO checkCOMMENT ownership
						r.Patch("/", app.checkPostOwnership("moderator", app.updateCommentHandler))
						r.Delete("/", app.checkPostOwnership("admin", app.deleteCommentHandler))
					})
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				// r.Use(app.AuthTokenMiddleware())
				//FIX: // r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)
				// I made these endpoint as PUT to show clients that these are idempotent
				// Meaning they will not create a resource if sent several times
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				// r.Use(app.AuthTokenMiddleware())
				r.Get("/feed", app.getUserFeedHandler)
			})

		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler) // create user in "users" with is_active=false and token in "user_invitations"
			// r.Post("/token", app.createTokenHandler) // update token of already created user with is_active=false
			r.Post("/login", app.loginUserHandler)
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {

	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 60, // time it takes for server to send response back to client
		ReadTimeout:  time.Second * 10, // time it takes for server from receiving the first packet to the last packet after tcp connection is formed
		IdleTimeout:  time.Minute,      // time it takes for tcp connection remains alive after sending response back to the client
	}

	app.logger.Infof("Launching server on port%v", app.config.addr)

	return srv.ListenAndServe()
}
