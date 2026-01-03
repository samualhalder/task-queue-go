package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/samualhalder/task-queue-go/internal/env"
)

func main(){

	err:=godotenv.Load()
	if err!=nil{
		fmt.Print("error while loading env")
	}

	addr:=env.String("ADDR","899")

	fmt.Print("hello world",addr)
}