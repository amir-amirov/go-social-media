package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amir-amirov/go-social-media/docs"
	"github.com/amir-amirov/go-social-media/internal/env"
	"github.com/amir-amirov/go-social-media/internal/mailer"
	"github.com/amir-amirov/go-social-media/internal/ratelimiter"
	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/amir-amirov/go-social-media/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config      config
	store       store.Storage
	logger      *zap.SugaredLogger
	mailer      mailer.Client
	cacheStore  cache.Storage
	rateLimiter ratelimiter.Limiter
}

type config struct {
	addr        string
	apiURL      string
	db          dbConfig
	env         string
	mail        mailConfig
	redis       redisConfig
	rateLimiter ratelimiter.Config
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

type redisConfig struct {
	addr     string
	password string
	db       int
	enabled  bool
}

func newApplication(cfg config, store store.Storage, logger *zap.SugaredLogger, mailer mailer.Client, cacheStore cache.Storage, rateLimiter ratelimiter.Limiter) *application {

	return &application{
		config:      cfg,
		store:       store,
		logger:      logger,
		mailer:      mailer,
		cacheStore:  cacheStore,
		rateLimiter: rateLimiter,
	}
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5173")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	if app.config.rateLimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	docURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	r.Route("/v1", func(r chi.Router) {

		r.Get("/health", app.healthCheckHandler)

		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/users", func(r chi.Router) {
				r.Get("/{userID}", app.getUserPostsHandler)
			})

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))

				r.Route("/comments", func(r chi.Router) {

					r.Post("/", app.createCommentHandler)

					r.Route("/{commentID}", func(r chi.Router) {
						r.Use(app.commentsContextMiddleware)

						r.Get("/", app.getCommentsHandler)
						r.Patch("/", app.checkPostOwnership("moderator", app.updateCommentHandler))
						r.Delete("/", app.checkPostOwnership("admin", app.deleteCommentHandler))
					})
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				// I made these endpoint as PUT to show clients that these are idempotent
				// Meaning they will not create a resource if sent several times
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
				r.Get("/follow-stats", app.getFollowStatsHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
				r.Get("/top", app.getTopUsersHandler)
			})

		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler) // create user in "users" with is_active=false and token in "user_invitations"
			r.Post("/token", app.resendTokenHandler) // update token of already created user with is_active=false
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

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infof("Launching server on port%v", app.config.addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
