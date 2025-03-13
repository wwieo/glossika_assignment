package gormDao

import (
	"glossika/service/internal/errorx"
	"glossika/service/internal/model"
	orm "glossika/service/repository"
	"gorm.io/gorm"
)

type accountDao struct {
	db *gorm.DB
}

func NewAccountDao(db *gorm.DB) orm.AccountDao {
	return &accountDao{
		db: db,
	}
}

func (dao *accountDao) Create(account *model.Account) error {
	err := dao.db.Table(model.TableAccount.String()).
		Create(account).Error
	if err != nil {
		if errorx.IsMySQLDuplicateEntry(err) {
			return errorx.RecordExisted
		}
		return err
	}
	return nil
}

func (dao *accountDao) Get(email string) (*model.Account, error) {
	var account *model.Account
	err := dao.db.Table(model.TableAccount.String()).
		Where("email = ?", email).
		First(&account).Error

	if err != nil {
		return nil, err
	}
	return account, nil
}

func (dao *accountDao) UpdateVerifyStatus(email string, status model.Status) error {
	if err := dao.db.Table(model.TableAccount.String()).
		Model(&model.Account{}).
		Where("email = ?", email).
		Update("status", status).
		Error; err != nil {
		return err
	}
	return nil
}
