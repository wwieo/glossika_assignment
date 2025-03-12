package orm

import "glossika/service/internal/model"

type AccountDao interface {
	Create(account *model.Account) error
	Get(email string) (*model.Account, error)
}
