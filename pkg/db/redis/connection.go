package rdb

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
	once        sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		initRedis()
	})
	return redisClient
}

func initRedis() {
	redisURI := os.Getenv("REDIS_URI")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisURI == "" || redisPassword == "" {
		panic("REDIS_URI and REDIS_PASSWORD environment variables must be set")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: redisPassword,
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}
	fmt.Println("Redis connected:", pong)
}
