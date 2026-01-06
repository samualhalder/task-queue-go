package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/samualhalder/task-queue-go/internal/store"
)

type application struct{
	config config
	store *store.Store
}

type config struct{
	addr string
	dbConfig dbConfig
	env string
}

type dbConfig struct{
	addr string
	maxOpenConn int
	maxIdlConn int
	maxIdlTime string
}


func(app *application) mount() http.Handler{
	r:=chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

    r.Route("/api/v1",func(r chi.Router){

	})

	return r
}

func(app *application) run(mux http.Handler) error {
	srv:=http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Second,
	}
	shutdown:=make(chan error)
	go func(){
		quit:=make(chan os.Signal)

		signal.Notify(quit, syscall.SIGINT,syscall.SIGTERM)
		s:= <-quit

		ctx,cancel:=context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()
		fmt.Printf("signal caught",s.String())
		shutdown <-srv.Shutdown(ctx)
	}()


	fmt.Print("server is running on port " , app.config.addr)
	err:=srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed){
		return err
	}
	err= <-shutdown

	if err!=nil{
		return err
	}
	return nil
}