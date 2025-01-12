package redis

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// Create a Redis client with the specified details
func CreateClient() (*redis.Client, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	address := os.Getenv("REDIS_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
		Protocol: 2,
	})

	return client, nil
}
