package main

import (
	"time"

	"github.com/amir-amirov/go-social-media/internal/db"
	"github.com/amir-amirov/go-social-media/internal/env"
	"github.com/amir-amirov/go-social-media/internal/mailer"
	"github.com/amir-amirov/go-social-media/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Go Social Media API
//	@description	API for Social Media written in Go
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api/v1

// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost:5431/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days to accept invitation
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendgrid: sendGridConfig{
				apiKey: env.GetString("SEND_GRID_API_KEY", ""),
			},
			mailtrap: mailTrapConfig{
				apiKey: env.GetString("MAIL_TRAP_API_KEY", ""),
			},
		},
	}

	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns)
	if err != nil {
		logger.Fatal("Unable to connect to database..")
	}

	defer db.Close()

	store := store.NewPostgresStorage(db)

	// mailer := mailer.NewSendGrid(cfg.mail.fromEmail, cfg.mail.sendgrid.apiKey)
	mailer, err := mailer.NewMailTrapClient(cfg.mail.mailtrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	app := newApplication(cfg, store, logger, mailer)

	mux := app.mount()

	if err := app.run(mux); err != nil {
		logger.Fatal("Unable to launch server..")
	}
}
