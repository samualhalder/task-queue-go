package worker_main

import (
	"context"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/samualhalder/task-queue-go/internal/app"
	"github.com/samualhalder/task-queue-go/internal/db"
	"github.com/samualhalder/task-queue-go/internal/env"
	"github.com/samualhalder/task-queue-go/internal/logger"
	"github.com/samualhalder/task-queue-go/internal/queue"
	"github.com/samualhalder/task-queue-go/internal/redis"
	"github.com/samualhalder/task-queue-go/internal/store"
	"github.com/samualhalder/task-queue-go/internal/worker"
)

func main(){
	err:=godotenv.Load()
	if err!=nil{
		panic("error while loading env in worker")
	}
	
	cnf:=app.Config{
		DBConfig: app.DBConfig{
			Addr: env.String("DB_ADDR",""),
			MaxOpenConn: env.Int("MAX_OPEN_CONN",30),
			MaxIdlConn: env.Int("MAX_IDLE_CONN",30),
			MaxIdlTime: env.String("MAX_IDLE_TIME","15m"),
		},
		Env:env.String("ENV","dev"),
		RedisCnf: app.RedisConfig{
			Addr: env.String("REDIS_ADDR",""),
			Password: env.String("REDIS_PASSWORD",""),
			DB: env.Int("REDIS_DB",0),
			Enabled: env.Bool("REDIS_ENABLED",true),
		},
	}

	db,err:=db.New(cnf.DBConfig.Addr,cnf.DBConfig.MaxOpenConn,cnf.DBConfig.MaxIdlConn,cnf.DBConfig.MaxIdlTime)

	if err!=nil{
		panic("error while db setup in worker: "+err.Error())
	}
	store:=store.NewStore(db)

	 logger,err:= logger.New(cnf.Env)
	 defer logger.Sync()
	 if err!=nil{
		panic("error logger setup : "+err.Error())
	 }

	 
	rdb:=redis.New(cnf.RedisCnf.Password,cnf.RedisCnf.Addr,cnf.RedisCnf.DB)
	taskQueue:= queue.NewRedisQueue(rdb,"queue:tasks:default")

	worker:= worker.NewWorker(0,taskQueue,store.Task,nil,logger.Sugar())

	ctx,stop:= signal.NotifyContext(context.Background(),os.Interrupt)
	defer stop()

	worker.Start(ctx)
}


