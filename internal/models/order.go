package models

import "time"

type Order struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Orden_compra string `gorm:"primaryKey;size:255;not null"`
	CreatedAt    time.Time
	Filename     string `gorm:"size:255; not null"`
}

func (Order) TableName() string {
	return "pedidos"
}
