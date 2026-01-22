package mailer

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
	taskerrors "github.com/samualhalder/task-queue-go/internal/errors"
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

func (r *resendMailer) Send(ctx context.Context, data EmailMessage) *taskerrors.TaskError {
	fmt.Print("tried to send mail")

	params := &resend.SendEmailRequest{
		From:    r.From,
		To:      data.To,
		Html:    data.HTML,
		Subject: data.Subject,
	}
	_, err := r.Client.Emails.Send(params)
	fmt.Printf("sending mail :", "mail", params.To)
	if err != nil {
		return taskerrors.Retryable(fmt.Errorf("error while sending mail: %w", err))
	}
	return nil
}
