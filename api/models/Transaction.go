package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID                   uint           `gorm:"primaryKey" json:"id" form:"id"`
	TransactionName      string         `gorm:"transaction_name" form:"transaction_name"`
	TransactionReceiver  string         `json:"transaction_receiver" form:"transaction_receiver"`
	TransactionDepositor string         `json:"transaction_depositor" form:"transaction_depositor"`
	TransactionAmount    int64          `json:"transaction_amount" form:"transaction_amount"`
	TransactionDate      time.Time      `time_format:"2006-01-02" json:"transaction_date" form:"transaction_date"`
	IsDebit              bool           `json:"is_debit" form:"is_debit"`
	CreatedAt            time.Time      `json:"-"`
	UpdatedAt            time.Time      `json:"-"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Transaction) CreateTransaction(db *gorm.DB) (*Transaction, error) {
	err := db.Create(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	return t, nil
}

func (t *Transaction) GetAllTransaction(db *gorm.DB) (*[]Transaction, error) {
	ts := []Transaction{}
	err := db.Model(&Account{}).Find(&ts).Error
	if err != nil {
		return &[]Transaction{}, err
	}
	return &ts, nil
}
