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
	// Prefer REDIS_URL (used by Railway and other cloud providers)
	if url := os.Getenv("REDIS_URL"); url != "" {
		opt, err := redis.ParseURL(url)
		if err != nil {
			log.Fatalf("Failed to parse REDIS_URL: %v", err)
		}
		initClient(opt)
		return
	}

	// Fallback to individual Redis settings (for local development)
	host := getenv("REDIS_HOST", "localhost")
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
	log.Printf("Attempting to connect to Redis at: %s", opt.Addr)
	redisClient = redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if ping, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Printf("Redis connection failed at %s: %v", opt.Addr, err)
		log.Fatalf("Failed to connect to Redis. Check REDIS_URL or Redis service configuration.")
	} else {
		log.Printf("Redis connection established: %s", ping)
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
