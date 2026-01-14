package queue

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)



type RedisQueue struct{
	rdb *redis.Client
	key string
}

func NewRedisQueue(rdb *redis.Client,key string) *RedisQueue{
	return &RedisQueue{
		rdb:rdb,
		key:key,
	}
}

func(r *RedisQueue) Push(ctx context.Context,uuid uuid.UUID) error{
	return r.rdb.RPush(ctx,r.key,uuid.String()).Err()
}

func(r *RedisQueue) Pop(ctx context.Context) (uuid.UUID,error){
	// fmt.Print("queue is poping ")
	res,err:= r.rdb.BLPop(ctx,0,r.key).Result()
	if err!=nil{
		return uuid.Nil,err
	}
	UUID,err:=uuid.Parse(res[1])
	if err!=nil{
		return uuid.Nil,err
	}
	return UUID,nil
}