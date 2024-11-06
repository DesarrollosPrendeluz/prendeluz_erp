package models

import "time"

// TODO add foreign keys
type ErrorOrder struct {
	ID        uint64 `gorm:"primary_key;not null"`
	Order     string `gorm:"size:255;not null;column:order_id"`
	Main_Sku  string `gorm:"size:255;not null"`
	Error     string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ErrorOrder) TableName() string {
	return "order_errors"
}
