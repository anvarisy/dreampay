package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	AccountMobile string         `gorm:"primaryKey; size 16" json:"account_mobile" form:"account_mobile"`
	AccountName   string         `json:"account_name" form:"account_name"`
	AccountStatus string         `json:"account_status" form:"account_status" gorm:"size:2"`
	IsActive      bool           `json:"is_active" form:"is_active" gorm:"default:false"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
type AccountRegister struct {
	AccountMobile string ` json:"account_mobile" form:"account_mobile"`
	AccountName   string `json:"account_name" form:"account_name"`
	AccountStatus string `json:"account_status" form:"account_status"`
}

type AccountVerification struct {
	AccountMobile string ` json:"account_mobile" form:"account_mobile"`
}

type AccountLogin struct {
	AccountMobile string ` json:"account_mobile" form:"account_mobile"`
}

func (a *Account) GetAllAccount(db *gorm.DB) (*[]Account, error) {
	account := []Account{}
	err := db.Model(&Account{}).Find(&account).Error
	if err != nil {
		return &[]Account{}, err
	}
	return &account, nil
}

func (a *Account) CreateAccount(db *gorm.DB) (*Account, error) {
	err := db.Create(&a).Error
	if err != nil {
		return &Account{}, err
	}
	return a, nil
}
