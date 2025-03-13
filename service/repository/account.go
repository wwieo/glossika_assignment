package orm

import "glossika/service/internal/model"

type AccountDao interface {
	Create(account *model.Account) error
	Get(email string) (*model.Account, error)
	UpdateVerifyStatus(email string, status model.Status) error
}
