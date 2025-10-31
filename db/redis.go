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
	redisClient = redis.NewClient(opt)

	// Retry settings (tunable via env)
	maxAttempts := 15 // ~60-70s total with backoff
	if v := os.Getenv("REDIS_MAX_ATTEMPTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxAttempts = n
		}
	}

	baseDelay := 1 * time.Second
	if v := os.Getenv("REDIS_RETRY_BASE_MS"); v != "" {
		if ms, err := strconv.Atoi(v); err == nil && ms > 0 {
			baseDelay = time.Duration(ms) * time.Millisecond
		}
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Printf("Attempting to connect to Redis at: %s (attempt %d/%d)", opt.Addr, attempt, maxAttempts)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if ping, err := redisClient.Ping(ctx).Result(); err == nil {
			cancel()
			log.Printf("Redis connection established: %s", ping)
			return
		} else {
			lastErr = err
			cancel()
			// Exponential backoff with jitter
			sleep := baseDelay * time.Duration(1<<uint(min(attempt-1, 5)))
			if sleep > 8*time.Second {
				sleep = 8 * time.Second
			}
			time.Sleep(sleep)
		}
	}
	log.Printf("Redis connection failed at %s: %v", opt.Addr, lastErr)
	log.Fatalf("Failed to connect to Redis after retries. Ensure REDIS_URL is set (Railway) or Redis is reachable.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
