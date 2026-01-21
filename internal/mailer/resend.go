package mailer

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
)

type resendMailer struct {
	From   string
	ApiKey string
	Client *resend.Client
}

func NewResendMailer(apiKey, from string) *resendMailer {
	client := resend.NewClient(apiKey)
	return &resendMailer{
		ApiKey: apiKey,
		Client: client,
		From:   from,
	}
}

func (r *resendMailer) Send(ctx context.Context, data EmailMessage) error {
	fmt.Print("tried to send mail")
	return fmt.Errorf("mail service is down")

	params := &resend.SendEmailRequest{
		From:    r.From,
		To:      data.To,
		Html:    data.HTML,
		Subject: data.Subject,
	}
	_, err := r.Client.Emails.Send(params)
	fmt.Printf("sending mail :", "mail", params.To)
	if err != nil {
		return fmt.Errorf("error while sending mail: %w", err)
	}
	return nil
}
