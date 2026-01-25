package app

import (
	"net/http"

	jsonresponse "github.com/samualhalder/task-queue-go/internal/json_response"
)

func (app *Application) HealthChecker(w http.ResponseWriter, r *http.Request) {

	jsonresponse.Success(w, http.StatusOK, "done", nil)
}