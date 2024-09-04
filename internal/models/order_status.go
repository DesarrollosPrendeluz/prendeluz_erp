package models

import "time"

type OrderStatus struct {
	ID        uint64 `gorm:"autoIncrement"`
	Name      string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (OrderStatus) TableName() string {
	return "order_statuses"
}
