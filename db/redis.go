package db

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func StartRedis() {
	redisURL := os.Getenv("REDIS_DB_URL")

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		fmt.Println("Error parsing Redis URL:", err)
		return
	}

	redisClient = redis.NewClient(opt)

	ctx := context.Background()
	ping, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	fmt.Println("Redis Ping:", ping)

	fmt.Println("Connected to Redis!")
}

func GetRedisClient() *redis.Client {
	return redisClient
}
