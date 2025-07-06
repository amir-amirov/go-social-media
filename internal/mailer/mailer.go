package mailer

import "embed"

const (
	FromName            = "GopherSocial"
	maxRetires          = 3
	UserWelcomeTemplate = "user_invitation_code.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}
