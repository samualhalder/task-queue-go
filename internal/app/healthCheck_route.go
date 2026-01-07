package app

import "github.com/go-chi/chi/v5"

func(app *Application) HealthCheckRoute(r chi.Router){
	r.Route("/health",func(r chi.Router){
		r.Get("/",app.HealthChecker)
	})
}