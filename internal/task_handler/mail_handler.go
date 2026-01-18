package taskhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/samualhalder/task-queue-go/internal/mailer"
	"github.com/samualhalder/task-queue-go/internal/model"
	"go.uber.org/zap"
)

type MailHandler struct{
  Mail mailer.Client
  Logger *zap.SugaredLogger
 // TODO : add logger 
}


func (m *MailHandler) Handle(ctx context.Context, task *model.Task) error {
	var dataPayload model.SendMailTemplate

	// 1. Decode payload

	if err := json.Unmarshal(task.Payload, &dataPayload); err != nil {
		return fmt.Errorf("invalid send-mail payload: %w", err)
	}

	// 2. Load template
	tmpl, err := template.ParseFS(
		mailer.FS,
		"templates/"+dataPayload.Template +".tmpl",
	)
	if err != nil {
		return fmt.Errorf("no such template present: %w", err)
	}

	// 3. Render template with data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, dataPayload.Data); err != nil {
		return fmt.Errorf("template render failed: %w", err)
	}

	// 4. Resolve subject (template default or user override)
	subject := dataPayload.Subject
	if subject == "" {
		subject = "Welcome" // or fetch from template metadata
	}

	// 5. Send email (mailer sees final content only)
	return m.Mail.Send(ctx, mailer.EmailMessage{
		To:      dataPayload.To,
		Subject: subject,
		HTML:    body.String(),
	})
}
