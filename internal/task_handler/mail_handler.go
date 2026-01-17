package taskhandler

import (
	"context"

	"github.com/samualhalder/task-queue-go/internal/mailer"
	"github.com/samualhalder/task-queue-go/internal/model"
	"go.uber.org/zap"
)

type MailHandler struct{
  mail mailer.Client
  logger *zap.SugaredLogger
 // TODO : add logger 
}


func(m *MailHandler) Handle(ctx context.Context, task *model.Task) error{
	// TODO: this will have a mailer that should have a send 
	// TODO: extract the payload and call send method
	return m.mail.Send("","",string(task.Payload),"",false) // TODO put full logic here
}