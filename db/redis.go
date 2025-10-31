package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func StartRedis() {
	// Prefer REDIS_URL; fallback to individual host/port/password/db vars
	if url := os.Getenv("REDIS_URL"); url != "" {
		opt, err := redis.ParseURL(url)
		if err != nil {
			log.Fatalf("Failed to parse REDIS_URL: %v", err)
		}
		initClient(opt)
		return
	}

	host := getenv("REDIS_HOST", "redis")
	port := getenv("REDIS_PORT", "6379")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := getenv("REDIS_DB", "0")
	dbIdx, err := strconv.Atoi(dbStr)
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value '%s': %v", dbStr, err)
	}

	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       dbIdx,
	}
	initClient(opt)
}

func initClient(opt *redis.Options) {
	redisClient = redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if ping, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", opt.Addr, err)
	} else {
		fmt.Println("Redis connection established:", ping)
	}
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
