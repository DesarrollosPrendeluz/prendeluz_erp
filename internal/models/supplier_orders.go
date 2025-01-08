package models

import "time"

type SupplierOrder struct {
	FatherOrderID uint64 `gorm:"not null"`
	SupplierID    uint64 `gorm:"not null"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

func (SupplierOrder) TableName() string {
	return "supplier_orders"
}
