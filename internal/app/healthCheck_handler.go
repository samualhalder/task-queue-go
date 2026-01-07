package app

import "net/http"

func(app *Application) HealthChecker(w  http.ResponseWriter,r *http.Request){
	
	w.Write([]byte("its good"))
}