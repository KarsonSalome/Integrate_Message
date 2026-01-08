package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	log.Println("Redis connected")
}
