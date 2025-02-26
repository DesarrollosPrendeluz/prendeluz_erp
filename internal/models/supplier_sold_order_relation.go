package models

import (
	"time"

	"gorm.io/gorm"
)

// SupplierSoldOrderRelation representa la relación entre supplier y sold_order
type SupplierSoldOrderRelation struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	SupplierID  uint `gorm:"column:supplier_father_order_id;not null;index"`
	SoldOrderID uint `gorm:"column:seller_father_order_id;not null;index"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

// Tabla personalizada (opcional)
func (SupplierSoldOrderRelation) TableName() string {
	return "supplier_sold_order_relation"
}

// Claves foráneas en la base de datos
func (relation *SupplierSoldOrderRelation) BeforeCreate(tx *gorm.DB) (err error) {
	// Aquí podrías agregar validaciones antes de insertar en la base de datos
	return nil
}
