package models

import (
	"time"

	"gorm.io/gorm"
)

type Withdraw struct {
	ID        uint      `gorm:"primaryKey"`
	Amount    int64     `json:"amount"`
	SellerID  string    `gorm:"size:16" json:"seller_id"`
	S         Account   `gorm:"foreignKey:SellerID; size:16"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (w *Withdraw) CreateWithdraw(db *gorm.DB) (*Withdraw, error) {
	err := db.Create(&w).Error
	if err != nil {
		return &Withdraw{}, err
	}
	return w, nil
}
