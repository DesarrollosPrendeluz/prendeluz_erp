package models

import "time"

type OrderStatus string

const (
	InProgress OrderStatus = "in progress"
	Received   OrderStatus = "received"
	Completed  OrderStatus = "completed"
	Canceled   OrderStatus = "canceled"
)

type Order struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Orden_compra string `gorm:"primaryKey;size:255;not null"`
	CreatedAt    time.Time
	Filename     string      `gorm:"size:255; not null"`
	Status       OrderStatus `gorm:"size:255;not null; type:enum('in progress','received','canceled','completed');default:'received'"`
}

func (Order) TableName() string {
	return "pedidos"
}
