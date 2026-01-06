package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/samualhalder/task-queue-go/internal/db"
	"github.com/samualhalder/task-queue-go/internal/env"
	"github.com/samualhalder/task-queue-go/internal/store"
)

func main(){

	err:=godotenv.Load()
	if err!=nil{
		fmt.Print("error while loading env")
	}
	cnf:=config{
		addr: env.String("ADDR",":8080"),
		dbConfig: dbConfig{
			addr: env.String("DB_ADDR",""),
			maxOpenConn: env.Int("MAX_OPEN_CONN",30),
			maxIdlConn: env.Int("MAX_IDLE_CONN",30),
			maxIdlTime: env.String("MAX_IDLE_TIME","15m"),
		},
		env:env.String("ENV","dev"),
	}

	db,err:=db.New(cnf.dbConfig.addr,cnf.dbConfig.maxIdlConn,cnf.dbConfig.maxIdlConn,cnf.dbConfig.maxIdlTime)

	if err!=nil{
		fmt.Printf("error while loading database",err.Error())
		return 
	}
	store:=store.NewStore(db)

	app:=application{
		config: cnf,
		store: store,
	}

	mux:=app.mount()
	fmt.Print("route setup completed")
	app.run(mux)
}