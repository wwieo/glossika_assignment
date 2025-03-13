package caches

import "context"

type AccountCache interface {
	SetVerifyCode(ctx context.Context, verifyCode string, email string) error
	GetEmailFromVerifyCode(ctx context.Context, verifyCode string) (string, error)
	DeleteVerifyCode(ctx context.Context, verifyCode string)
}
