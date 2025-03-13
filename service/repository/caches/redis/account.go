package redisDao

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"glossika/service/repository/caches"
	"time"
)

func NewAccountCache(r *redis.Client) caches.AccountCache {
	return &accountCache{
		prefixVerifyCodeKey: "verify_code",
		client:              r,
	}
}

type accountCache struct {
	prefixVerifyCodeKey string
	client              *redis.Client
}

func (dao *accountCache) buildInfoKey(verifyCode string) string {
	return fmt.Sprintf("%s:%s", dao.prefixVerifyCodeKey, verifyCode)
}

func (dao *accountCache) SetVerifyCode(ctx context.Context, verifyCode string, email string) error {
	return dao.client.Set(ctx, dao.buildInfoKey(verifyCode), email, 10*time.Minute).Err()
}

func (dao *accountCache) GetEmailFromVerifyCode(ctx context.Context, verifyCode string) (string, error) {
	email, err := dao.client.Get(ctx, dao.buildInfoKey(verifyCode)).Result()
	if err != nil {
		return "", err
	}
	return email, nil
}

func (dao *accountCache) DeleteVerifyCode(ctx context.Context, verifyCode string) {
	dao.client.Del(ctx, dao.buildInfoKey(verifyCode))
}
