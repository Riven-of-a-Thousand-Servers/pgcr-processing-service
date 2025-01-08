package redis

import (
	"context"
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	dto "rivenbot/types/dto"
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

// Returns a given manifest entity based on a hash
func GetManifestEntity(client *redis.Client, hash string) (*dto.ManifestObject, error) {
	ctx := context.Background()
	result, err := client.Get(ctx, hash).Result()
	if err != nil {
		return nil, err
	}

	var response *dto.ManifestObject
	if err := json.Unmarshal([]byte(result), response); err != nil {
		return nil, err
	}

	return response, nil
}
