package accountCtrl

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"glossika/service/internal/errorx"
	"glossika/service/internal/model"
	boAccount "glossika/service/internal/model/bo/account"
	"glossika/service/internal/utils"
	"glossika/service/repository/orm/gormDao"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type AccountCtrl interface {
	Register(ctx context.Context, args boAccount.RegisterArgs) (boAccount.RegisterReply, error)
	Login(ctx context.Context, args boAccount.LoginArgs) (boAccount.LoginReply, error)
}

func New(pack accountCtrlPack) AccountCtrl {
	return &accountCtrl{
		pack: pack,
	}
}

type accountCtrlPack struct {
	dig.In

	MySQLGlossika *gorm.DB `name:"glossika"`
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
	dao := gormDao.NewAccountDao(ctrl.pack.MySQLGlossika)
	// TODO verify logic
	err := dao.Create(account)
	if err != nil {
		return boAccount.RegisterReply{}, err
	}
	return boAccount.RegisterReply{}, nil
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
	return reply, nil
}
