package models

import "time"

// TODO add foreign keys
type ErrorOrder struct {
	ID        uint64 `gorm:"primary_key;not null"`
	Order     string `gorm:"size:255;not null;column:orden_pedido"`
	Main_Sku  string `gorm:"size:255;not null"`
	Error     string `gorm:"size:255;not null"`
	CreatedAt time.Time
}

func (ErrorOrder) TableName() string {
	return "pedido_error"
}
