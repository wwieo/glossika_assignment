package accountCtrl

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"glossika/service/internal/config"
	"glossika/service/internal/errorx"
	"glossika/service/internal/model"
	boAccount "glossika/service/internal/model/bo/account"
	"glossika/service/internal/utils"
	redisDao "glossika/service/repository/caches/redis"
	"glossika/service/repository/orm/gormDao"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type AccountCtrl interface {
	Register(ctx context.Context, args boAccount.RegisterArgs) (boAccount.RegisterReply, error)
	Login(ctx context.Context, args boAccount.LoginArgs) (boAccount.LoginReply, error)
	Verify(ctx context.Context, args boAccount.VerifyArgs) (boAccount.VerifyReply, error)
}

func New(pack accountCtrlPack) AccountCtrl {
	return &accountCtrl{
		pack: pack,
	}
}

type accountCtrlPack struct {
	dig.In

	ServiceAddress config.ServiceAddress
	MySQLGlossika  *gorm.DB      `name:"glossika"`
	RedisGlossika  *redis.Client `name:"glossika"`
}

type accountCtrl struct {
	pack accountCtrlPack
}

func passwordHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (ctrl *accountCtrl) Register(ctx context.Context, args boAccount.RegisterArgs) (boAccount.RegisterReply, error) {
	account := &model.Account{
		ID:       utils.GetSnowflakeIDInt64(),
		Email:    args.Email,
		Status:   model.StatusUnVerified,
		Password: passwordHash(args.Password),
	}
	var reply boAccount.RegisterReply
	dao := gormDao.NewAccountDao(ctrl.pack.MySQLGlossika)
	err := dao.Create(account)
	if err != nil {
		return reply, err
	}

	rao := redisDao.NewAccountCache(ctrl.pack.RedisGlossika)
	verifyCode := utils.GenerateRandomString(15)
	err = rao.SetVerifyCode(ctx, verifyCode, args.Email)
	if err != nil {
		return reply, err
	}

	// just for temp hard code usage
	reply.VerifyURL = fmt.Sprintf("localhost%s/glossika/account/verify/%s",
		ctrl.pack.ServiceAddress.Glossika, verifyCode)
	return reply, nil
}

func (ctrl *accountCtrl) Login(ctx context.Context, args boAccount.LoginArgs) (boAccount.LoginReply, error) {
	dao := gormDao.NewAccountDao(ctrl.pack.MySQLGlossika)
	reply := boAccount.LoginReply{}
	account, err := dao.Get(args.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reply, errorx.NoRecord
		}
		return reply, err
	}

	if passwordHash(args.Password) != account.Password {
		return reply, errorx.PasswordIncorrect
	}
	reply.IsVerified = account.Status == model.StatusVerified
	return reply, nil
}

func (ctrl *accountCtrl) Verify(ctx context.Context, args boAccount.VerifyArgs) (boAccount.VerifyReply, error) {
	var reply boAccount.VerifyReply
	rao := redisDao.NewAccountCache(ctrl.pack.RedisGlossika)
	email, err := rao.GetEmailFromVerifyCode(ctx, args.VerifyCode)
	if err != nil {
		return reply, errorx.InvalidVerifyCode
	}

	dao := gormDao.NewAccountDao(ctrl.pack.MySQLGlossika)
	err = dao.UpdateVerifyStatus(email, model.StatusVerified)
	if err != nil {
		return reply, err
	}

	rao.DeleteVerifyCode(ctx, args.VerifyCode)
	return reply, nil
}
