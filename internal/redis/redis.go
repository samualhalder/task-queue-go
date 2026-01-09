package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)



func New(password,addr string,db int) *redis.Client{
	rdb:= redis.NewClient(&redis.Options{
		Password: password,
		Addr: addr,
		DB: db,
	})
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second * 2)
	defer cancel()
	
	if err:= rdb.Ping(ctx).Err();err!=nil{
		panic("failed to connect with redis service: "+ err.Error())
	}
	return rdb
}