package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/samualhalder/task-queue-go/internal/dto"
	jsonresponse "github.com/samualhalder/task-queue-go/internal/json_response"
)

type taskCreatePayload struct {
	TaskName     string          `json:"task_name" validate:"required,min=3"`
	Payload      json.RawMessage `json:"payload" validate:"required,json"`
	ScheduledAt  *time.Time      `json:"scheduled_at" validate:"omitempty,gt"`
	Max_Attempts int             `json:"max_attempts" validate:"omitempty,min=1"`
}

func (app *Application) CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskData taskCreatePayload
	if err := jsonresponse.ReadJson(w, r, &taskData); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := jsonresponse.Validator.Struct(taskData); err != nil {
		app.badRequest(w, r, err)
		return
	}

	task := &dto.TaskResponse{
		TaskName:    taskData.TaskName,
		Payload:     taskData.Payload,
		ScheduledAt: taskData.ScheduledAt,
	}
	if taskData.Max_Attempts > 0 {
		task.Max_Attempts = min(taskData.Max_Attempts, app.Config.DefaultMaxAttempts)
	} else {
		task.Max_Attempts = app.Config.DefaultMaxAttempts
	}
	ctx := r.Context()

	if err := app.Store.Task.Create(ctx, task); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.Queue.Signal(ctx); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonresponse.Success(w, http.StatusCreated, "task created successfully", task); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
