package main

import (
	"github.com/joho/godotenv"
	"github.com/samualhalder/task-queue-go/internal/app"
	"github.com/samualhalder/task-queue-go/internal/db"
	"github.com/samualhalder/task-queue-go/internal/env"
	"github.com/samualhalder/task-queue-go/internal/logger"
	"github.com/samualhalder/task-queue-go/internal/queue"
	"github.com/samualhalder/task-queue-go/internal/redis"
	"github.com/samualhalder/task-queue-go/internal/store"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("error while loading env")
	}

	cnf := app.Config{
		Addr: env.String("ADDR", ":8080"),
		DBConfig: app.DBConfig{
			Addr:        env.String("DB_ADDR", ""),
			MaxOpenConn: env.Int("MAX_OPEN_CONN", 30),
			MaxIdlConn:  env.Int("MAX_IDLE_CONN", 30),
			MaxIdlTime:  env.String("MAX_IDLE_TIME", "15m"),
		},
		Env: env.String("ENV", "dev"),
		RedisCnf: app.RedisConfig{
			Addr:     env.String("REDIS_ADDR", ""),
			Password: env.String("REDIS_PASSWORD", ""),
			DB:       env.Int("REDIS_DB", 0),
			Enabled:  env.Bool("REDIS_ENABLED", true),
		},
		DefaultMaxAttempts: env.Int("MAX_ATTEMPTS", 5),
	}

	db, err := db.New(cnf.DBConfig.Addr, cnf.DBConfig.MaxOpenConn, cnf.DBConfig.MaxIdlConn, cnf.DBConfig.MaxIdlTime)

	if err != nil {
		panic("error while db setup: " + err.Error())
	}
	store := store.NewStore(db)

	logger, err := logger.New(cnf.Env)
	defer logger.Sync()
	if err != nil {
		panic("error logger setup : " + err.Error())
	}

	rdb := redis.New(cnf.RedisCnf.Password, cnf.RedisCnf.Addr, cnf.RedisCnf.DB)
	taskQueue := queue.NewRedisQueue(rdb, "queue:tasks:default")

	app := app.New(cnf, store, logger.Sugar(), taskQueue)

	if err := app.Run(); err != nil {
		panic(err)
	}

}
