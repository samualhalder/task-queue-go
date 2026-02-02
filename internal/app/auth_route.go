package app

import "github.com/go-chi/chi/v5"

func (app *Application) authRoute(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Route("/admin", func(r chi.Router) {
			r.Post("/get-token", app.GetToken)
		})
	})
}
