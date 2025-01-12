package redis

import (
	"context"
	"rivenbot/internal/dto"
)

type RedisClient interface {
	GetManifestEntity(ctx context.Context, hash string) (*dto.ManifestObject, error)
}
