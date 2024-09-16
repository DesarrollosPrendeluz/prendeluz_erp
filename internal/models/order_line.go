package models

import "time"

type OrderLine struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement"`
	OrderID       uint64 `gorm:"column:id_order;primaryKey;not null"`
	ItemId        string `gorm:"column:id_item;primaryKey;not null"`
	Amount        int64  `gorm:"column:quantity"`
	AmountRecived int64  `gorm:"column:recived_quantity"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

func (OrderLine) TableName() string {
	return "stock_deficits"
}
