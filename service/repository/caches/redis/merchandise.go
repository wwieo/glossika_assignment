package redisDao

import (
	"context"
	"github.com/redis/go-redis/v9"
	"glossika/service/repository/caches"
	"time"
)

func NewMerchandiseCache(r *redis.Client) caches.MerchandiseCache {
	return &merchandiseCache{
		prefixRecommendationKey: "merchandise:recommendation",
		client:                  r,
	}
}

type merchandiseCache struct {
	prefixRecommendationKey string
	client                  *redis.Client
}

func (dao *merchandiseCache) buildRecommendationKey() string {
	return dao.prefixRecommendationKey
}

func (dao *merchandiseCache) GetRecommendation(ctx context.Context) (string, error) {
	return dao.client.Get(ctx, dao.buildRecommendationKey()).Result()
}

func (dao *merchandiseCache) SetRecommendation(ctx context.Context, recommendation string) error {
	return dao.client.Set(ctx, dao.buildRecommendationKey(), recommendation, 10*time.Minute).Err()
}
