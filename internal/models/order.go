package models

import (
	"time"

	"gorm.io/gorm"
)

// type OrderStatus string

// const (
// 	InProgress OrderStatus = "in progress"
// 	Received   OrderStatus = "received"
// 	Completed  OrderStatus = "completed"
// 	Canceled   OrderStatus = "canceled"
// )

type Order struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement"`
	FatherOrderID uint64
	OrderStatusID uint64
	Code          string `gorm:"primaryKey;size:255;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	OrderStatus     OrderStatus `gorm:"foreignKey:ID;references:OrderStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	FatherOrder     FatherOrder `gorm:"foreignKey:ID;references:FatherOrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderLines      []OrderItem `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity        int64       `gorm:"-"`
	RecivedQuantity int64       `gorm:"-"`
	// Orden_compra string `gorm:"primaryKey;size:255;not null"`
	// Filename     string      `gorm:"size:255; not null"`
	// Status       OrderStatus `gorm:"size:255;not null; type:enum('in progress','received','canceled','completed');default:'received'"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderTotals struct {
	Total   int64
	Partial int64
}

func (o *Order) AfterFind(tx *gorm.DB) (err error) {
	var totals OrderTotals
	err = tx.Table("order_lines").
		Select("sum(quantity) AS total, sum(recived_quantity) AS partial").
		Where("order_id = ?", o.ID).
		Group("order_id").
		Order("order_id").
		Take(&totals).Error

	if err != nil {
		return err
	}
	o.Quantity = totals.Total
	o.RecivedQuantity = totals.Partial
	return nil
}
