package mailer

import (
	"context"
	"embed"
)

const (
	FromName                 = "GO-SOCIAL"
	MaxRetries               = 3
	UserRegisterMailTemplate = "registermail.tmpl"
)

//go:embed "templates"
var FS embed.FS

type EmailMessage struct{
	To []string
	Subject string
	HTML string
}

type Client interface {
	Send(context.Context, EmailMessage) error
}
