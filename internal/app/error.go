package app

import (
	"net/http"

	jsonresponse "github.com/samualhalder/task-queue-go/internal/json_response"
)

func (app *Application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Bad Request: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.Logger.Errorw("Bad Request", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	jsonresponse.Error(w, http.StatusBadRequest, err.Error())
}
func (app *Application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Not Found: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.Logger.Errorw("Not Found", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	jsonresponse.Error(w, http.StatusNotFound, "Record not found")
}
func (app *Application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Internal Server Error: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.Logger.Errorw("Internal Server Error", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	jsonresponse.Error(w, http.StatusBadRequest, "Something went wrong")
}
func (app *Application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Conflict Error: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.Logger.Errorw("Conflict Error", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	jsonresponse.Error(w, http.StatusBadRequest, err.Error())
}
func (app *Application) authorizationError(w http.ResponseWriter, r *http.Request, err error) {
	// log.Printf("Conflict Error: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.Logger.Errorw("Authorization Error", "Methode", r.Method, "Path", r.URL.Path, "error", err.Error())
	jsonresponse.Error(w, http.StatusUnauthorized, err.Error())
}

func (app *Application) forbiddenError(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Conflict Error: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.Logger.Warnw("Forbidden Error", "Methode", r.Method, "Path", r.URL.Path, "error", "forbidden")
	jsonresponse.Error(w, http.StatusForbidden, "forbidden")
}
func (app *Application) rateLimitExcedError(w http.ResponseWriter, r *http.Request, time string) {
	// log.Printf("Conflict Error: %s,err: %s,path: %s", r.Method, r.URL.Path, err.Error())
	app.Logger.Warnw("Rate Limit Exced Error", "Methode", r.Method, "Path", r.URL.Path, "retry after: "+time)
	w.Header().Set("Retry-After", time)
	jsonresponse.Error(w, http.StatusTooManyRequests, "rate limit exced")
}
