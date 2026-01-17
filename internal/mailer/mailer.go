package mailer

const (
	FromName                 = "GO-SOCIAL"
	MaxRetries               = 3
	UserRegisterMailTemplate = "registermail.tmpl"
)



type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
