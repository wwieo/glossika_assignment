package caches

import "context"

type MerchandiseCache interface {
	GetRecommendation(ctx context.Context) (string, error)
	SetRecommendation(ctx context.Context, recommendation string) error
}
