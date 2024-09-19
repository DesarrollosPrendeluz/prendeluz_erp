package models

import "time"

// TODO add foreign keys
type InOrderRelation struct {
	ID          uint64 `gorm:"primary_key;not null"`
	SupplierID  uint64 `gorm:"column:supplier_id;not null"`
	OrderLineID uint64 `gorm:"column:order_line_id;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (InOrderRelation) TableName() string {
	return "in_order_relations"
}
