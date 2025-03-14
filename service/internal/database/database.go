package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"glossika/service/internal/config"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type GlossikaOut struct {
	dig.Out

	MySQLGlossika *gorm.DB      `name:"glossika"`
	RedisGlossika *redis.Client `name:"glossika"`
}

const (
	mysqlGlossika = "glossika"
	redisGlossika = "glossika"
)

func NewGlossika(ctx context.Context, dbms config.DatabaseManagementSystem) GlossikaOut {
	return GlossikaOut{
		MySQLGlossika: newMySQL(mysqlGlossika, dbms.MySQLSystems[mysqlGlossika]),
		RedisGlossika: newRedis(ctx, redisGlossika, dbms.RedisSystems[redisGlossika]),
	}
}
