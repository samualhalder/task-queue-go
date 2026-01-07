package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/samualhalder/task-queue-go/internal/app"
	"github.com/samualhalder/task-queue-go/internal/db"
	"github.com/samualhalder/task-queue-go/internal/env"
	"github.com/samualhalder/task-queue-go/internal/store"
)

func main(){

	err:=godotenv.Load()
	if err!=nil{
		fmt.Print("error while loading env")
	}
	
	cnf:=app.Config{
		Addr: env.String("ADDR",":8080"),
		DBConfig: app.DBConfig{
			Addr: env.String("DB_ADDR",""),
			MaxOpenConn: env.Int("MAX_OPEN_CONN",30),
			MaxIdlConn: env.Int("MAX_IDLE_CONN",30),
			MaxIdlTime: env.String("MAX_IDLE_TIME","15m"),
		},
		Env:env.String("ENV","dev"),
	}

	db,err:=db.New(cnf.DBConfig.Addr,cnf.DBConfig.MaxOpenConn,cnf.DBConfig.MaxIdlConn,cnf.DBConfig.MaxIdlTime)

	if err!=nil{
		fmt.Printf("error while loading database",err.Error())
		return 
	}
	store:=store.NewStore(db)

	app:=app.New(cnf,store)

	if err:=app.Run();err!=nil{
		fmt.Print("error while running app",err)
	}
}