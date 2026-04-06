package repository

import (
	"cafe/domain"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisBlacklistRepository struct {
	redisClient *redis.Client
	prefix      string
}

func NewRedisBlacklistRepository(redisClient *redis.Client, prefix string) domain.BlacklistService {
	return &redisBlacklistRepository{
		redisClient: redisClient,
		prefix:      prefix,
	}
}

func (r *redisBlacklistRepository) Add(c context.Context, tokenID string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", r.prefix, tokenID)
	return r.redisClient.Set(c, key, "1", expiration).Err()
}

func (r *redisBlacklistRepository) AddBatch(c context.Context, tokenIDs []string, expiration time.Duration) error {
	pipe := r.redisClient.Pipeline()
	for _, tokenID := range tokenIDs {
		key := fmt.Sprintf("%s%s", r.prefix, tokenID)
		pipe.Set(c, key, "1", expiration)
	}
	_, err := pipe.Exec(c)
	return err
}

func (r *redisBlacklistRepository) IsBlacklisted(c context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("%s%s", r.prefix, tokenID)
	val, err := r.redisClient.Exists(c, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}
