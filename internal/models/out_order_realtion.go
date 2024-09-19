package models

import "time"

// TODO add foreign keys
type OutOrderRelation struct {
	ID          uint64 `gorm:"primary_key;not null"`
	ClientID    uint64 `gorm:"column:customer_id;not null"`
	OrderLineID uint64 `gorm:"column:order_line_id;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (OutOrderRelation) TableName() string {
	return "out_order_relations"
}
