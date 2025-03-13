package merchandiseCtrl

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	boMerchandise "glossika/service/internal/model/bo/merchandise"
	redisDao "glossika/service/repository/caches/redis"
	"go.uber.org/dig"
)

type MerchandiseCtrl interface {
	GetRecommendation(ctx context.Context, args boMerchandise.GetRecommendationArgs) (boMerchandise.GetRecommendationReply, error)
}

func New(pack merchandiseCtrlPack) MerchandiseCtrl {
	return &merchandiseCtrl{
		pack: pack,
	}
}

type merchandiseCtrlPack struct {
	dig.In

	RedisGlossika *redis.Client `name:"glossika"`
}

type merchandiseCtrl struct {
	pack merchandiseCtrlPack
}

func (ctrl *merchandiseCtrl) GetRecommendation(ctx context.Context, args boMerchandise.GetRecommendationArgs) (boMerchandise.GetRecommendationReply, error) {
	reply := boMerchandise.GetRecommendationReply{}
	rao := redisDao.NewMerchandiseCache(ctrl.pack.RedisGlossika)
	data, err := rao.GetRecommendation(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return reply, err
	}
	// assume redis has no recommendation, and get recommendation data from mysql
	if errors.Is(err, redis.Nil) {
		data = "mock_recommendation_data"
	}
	rao.SetRecommendation(ctx, data)
	reply.Recommendation = data
	return reply, nil
}
