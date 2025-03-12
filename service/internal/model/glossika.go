package model

import (
	"time"

	"gorm.io/gorm"
)

type Status int

const (
	StatusUnVerified Status = iota
	StatusVerified

	TableAccount TableName = "account"
)

type Account struct {
	ID        int64     `json:"-" gorm:"primaryKey;column:id;autoIncrement:false"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null;column:email;uniqueIndex"`
	Status    Status    `json:"status" gorm:"type:tinyint;not null;column:status"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null;column:password"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

func (*Account) TableName() string {
	return TableAccount.String()
}

func GlossikaMigrate(db *gorm.DB) error {
	if db.Migrator().HasTable(&Account{}) {
		return nil
	}
	return db.AutoMigrate(&Account{})
}
