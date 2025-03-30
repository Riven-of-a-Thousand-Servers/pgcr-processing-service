package redis

import (
	"context"
	"encoding/json"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	"github.com/redis/go-redis/v9"
)

type ManifestClient interface {
	GetManifestEntity(ctx context.Context, hash string) (*types.ManifestObject, error)
}

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(url string) *RedisService {
	redis := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	return &RedisService{Client: redis}
}

// Returns a given manifest entity based on a hash
func (r *RedisService) GetManifestEntity(ctx context.Context, hash string) (*types.ManifestObject, error) {
	result, err := r.Client.Get(ctx, hash).Result()
	if err != nil {
		return nil, err
	}

	var response *types.ManifestObject
	if err := json.Unmarshal([]byte(result), response); err != nil {
		return nil, err
	}

	return response, nil
}
