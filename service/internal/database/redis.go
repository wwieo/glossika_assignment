package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"glossika/service/internal/config"
	"log"
)

func newRedis(ctx context.Context, name string, r config.Redis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: r.Password,
		DB:       r.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to ping Redis: %s, err: %v", name, err)
	}

	log.Printf("Pinged successfully redis: %s", name)
	return client
}
