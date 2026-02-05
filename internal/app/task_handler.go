package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/samualhalder/task-queue-go/internal/dto"
	jsonresponse "github.com/samualhalder/task-queue-go/internal/json_response"
	metrics "github.com/samualhalder/task-queue-go/internal/matrics"
	"github.com/samualhalder/task-queue-go/internal/payloads"
)

type taskCreatePayload struct {
	TaskName     string          `json:"task_name" validate:"required,min=3"`
	Payload      json.RawMessage `json:"payload" validate:"required,json"`
	ScheduledAt  *time.Time      `json:"scheduled_at" validate:"omitempty,gt"`
	Max_Attempts int             `json:"max_attempts" validate:"omitempty,min=1"`
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
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
	var scheduledTime time.Time

	if taskData.ScheduledAt != nil {
		scheduledTime = *taskData.ScheduledAt
	} else {
		// default: run immediately
		scheduledTime = time.Now().UTC()
	}
	task := &dto.TaskResponse{
		TaskName:    taskData.TaskName,
		Payload:     taskData.Payload,
		ScheduledAt: &scheduledTime,
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
	metrics.TasksPending.Inc()
	if err := jsonresponse.Success(w, http.StatusCreated, "task created successfully", task); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *Application) ListTasks(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query()

	taskName := p.Get("taskName")
	status := p.Get("status")
	createdBefore := p.Get("createdBefore")
	createdAfter := p.Get("createdAfter")
	orderBy := p.Get("orderBy")
	sortOrder := p.Get("sortOrder")
	limit := p.Get("limit")
	skip := p.Get("skip")

	var createdBeforeTime *time.Time
	if createdBefore != "" {
		t, err := time.Parse(time.RFC3339, createdBefore)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
		createdBeforeTime = &t
	}

	var createdAfterTime *time.Time
	if createdAfter != "" {
		t, err := time.Parse(time.RFC3339, createdAfter)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
		createdAfterTime = &t
	}

	limitInt := 50
	if limit != "" {
		num, err := strconv.Atoi(limit)
		if err != nil || num <= 0 || num > 200 {
			app.badRequest(w, r, fmt.Errorf("invalid limit"))
			return
		}
		limitInt = num
	}

	skipInt := 0
	if skip != "" {
		num, err := strconv.Atoi(skip)
		if err != nil || num < 0 {
			app.badRequest(w, r, fmt.Errorf("invalid skip"))
			return
		}
		skipInt = num
	}

	if orderBy == "" {
		orderBy = "created_at"
	}

	allowedOrderBy := map[string]bool{
		"created_at": true,
	}

	if !allowedOrderBy[orderBy] {
		app.badRequest(w, r, fmt.Errorf("invalid orderBy"))
		return
	}

	if sortOrder == "" {
		sortOrder = "desc"
	}

	sortOrder = strings.ToLower(sortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		app.badRequest(w, r, fmt.Errorf("invalid sortOrder"))
		return
	}

	params := payloads.TaskQueryPayload{
		TaskName:      strPtr(taskName),
		Status:        strPtr(status),
		CreatedBefore: createdBeforeTime,
		CreatedAfter:  createdAfterTime,
		OrderBy:       orderBy,
		SortOrder:     sortOrder,
		Limit:         limitInt,
		Skip:          skipInt,
	}

	ctx := r.Context()
	tasks, err := app.Store.Task.GetTasks(ctx, params)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	fmt.Print("came jere ", len(tasks))

	if err := jsonresponse.Success(w, http.StatusOK, "", tasks); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// TODO : task store to fetch data
}

func (app *Application) GetTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	UUID, err := uuid.Parse(id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	ctx := r.Context()
	task, err := app.Store.Task.GetById(ctx, UUID)
	if err := jsonresponse.Success(w, http.StatusOK, "", task); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) RunTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := uuid.Parse(id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	err = app.Store.Task.MakeTaskEligible(r.Context(), id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := jsonresponse.Success(w, http.StatusOK, "Task executed succesfully", nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
func (app *Application) FailedTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := uuid.Parse(id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	err = app.Store.Task.MakeTaskFailed(r.Context(), id)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := jsonresponse.Success(w, http.StatusOK, "Task failed succesfully", nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
