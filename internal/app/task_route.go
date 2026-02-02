package app

import "github.com/go-chi/chi/v5"

func (app *Application) taskRoute(r chi.Router) {
	r.Route("/task", func(r chi.Router) {
		r.Post("/", app.CreateTask)
		r.Route("/admin", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/", app.ListTasks)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", app.GetTaskById)
				r.Post("/run", app.RunTask)
				r.Post("/cancel", app.FailedTask)

			})

		})
	})
}
