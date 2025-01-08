package models

import "time"

type SupplierOrder struct {
	FatherOrderID uint64 `gorm:"not null"`
	SupplierID    uint64 `gorm:"not null"`
	CreatedAt     *time.Time
	UpdatedAt     *time.Time

	Supplier *Supplier `gorm:"foreignKey:ID;references:SupplierID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (SupplierOrder) TableName() string {
	return "supplier_orders"
}
