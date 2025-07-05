package cache

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
	redisURL := getEnv("REDIS_URL", "localhost:6379")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
