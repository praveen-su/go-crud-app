package main

import (
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}
