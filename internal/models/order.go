package models

import "time"

// type OrderStatus string

// const (
// 	InProgress OrderStatus = "in progress"
// 	Received   OrderStatus = "received"
// 	Completed  OrderStatus = "completed"
// 	Canceled   OrderStatus = "canceled"
// )

type Order struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement"`
	OrderStatusID uint64
	OrderTypeID   uint64
	Code          string `gorm:"primaryKey;size:255;not null"`
	Filename      string `gorm:"column:file_name;size:255;notnull"`

	CreatedAt time.Time
	UpdatedAt time.Time

	OrderStatus OrderStatus `gorm:"foreignKey:OrderStatusID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderType   OrderType   `gorm:"foreignKey:OrderTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Orden_compra string `gorm:"primaryKey;size:255;not null"`
	// Filename     string      `gorm:"size:255; not null"`
	// Status       OrderStatus `gorm:"size:255;not null; type:enum('in progress','received','canceled','completed');default:'received'"`
}

func (Order) TableName() string {
	return "orders"
}
