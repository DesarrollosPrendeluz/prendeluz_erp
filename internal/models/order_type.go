package models

import "time"

type OrderType struct {
	ID          uint64 `gorm:"autoIncrement"`
	Name        string `gorm:"size:255;not null"`
	StockModify string `gorm:"size:255;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (OrderType) TableName() string {
	return "order_types"
}
