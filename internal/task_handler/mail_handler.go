package taskhandler

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/mailer"
	"github.com/samualhalder/task-queue-go/internal/model"
	"go.uber.org/zap"
)

type MailHandler struct{
  Mail mailer.Client
  Logger *zap.SugaredLogger
 // TODO : add logger 
}


func(m *MailHandler) Handle(ctx context.Context, task *model.Task) error{
	// TODO: this will have a mailer that should have a send 
	// TODO: extract the payload and call send method
	return m.Mail.Send("","",string(task.Payload),"",false) // TODO put full logic here
}