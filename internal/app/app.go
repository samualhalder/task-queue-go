package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/samualhalder/task-queue-go/internal/queue"
	"github.com/samualhalder/task-queue-go/internal/store"
	"go.uber.org/zap"
)

type Application struct {
	Config Config
	Store  *store.Store
	router http.Handler
	Logger *zap.SugaredLogger
	Queue  queue.Queue
}

type Config struct {
	MaxAttempts  int
	Addr         string
	DBConfig     DBConfig
	Env          string
	WorkersCount int
	RedisCnf     RedisConfig
	EmailCnf     EmailConfig
}

type DBConfig struct {
	Addr        string
	MaxOpenConn int
	MaxIdlConn  int
	MaxIdlTime  string
}

type RedisConfig struct {
	Password string
	Addr     string
	DB       int
	Enabled  bool
}

type EmailConfig struct {
	ApiKey string
	From   string
}

func New(cnf Config, store *store.Store, logger *zap.SugaredLogger, queue queue.Queue) *Application {
	app := &Application{
		Config: cnf,
		Store:  store,
		Logger: logger,
		Queue:  queue,
	}
	app.router = app.mount()
	app.Logger.Info("Route is mounted")
	return app
}

func (app *Application) mount() http.Handler {
	r := chi.NewRouter()

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

	r.Route("/api/v1", func(r chi.Router) {
		app.HealthCheckRoute(r)
		app.taskRoute(r)
	})

	return r
}

func (app *Application) Run() error {
	srv := http.Server{
		Addr:         app.Config.Addr,
		Handler:      app.router,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second,
	}
	shutdown := make(chan error, 1)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	app.Logger.Info("Server is running on port ", app.Config.Addr)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdown

	if err != nil {
		return err
	}
	return nil
}
