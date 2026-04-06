package bootstrap

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(env *Env) *redis.Client {
	redisHost := env.RedisHost
	redisUsername := env.RedisUsername
	redisPass := env.RedisPass

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Username: redisUsername,
		Password: redisPass,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Warn: Could not connect to Redis, Blacklist feature might not work: ", err)
	} else {
		log.Println("Successfully connected to Redis")
	}

	return rdb
}
