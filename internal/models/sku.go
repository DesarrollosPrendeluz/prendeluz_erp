package models

import "gorm.io/gorm"

type Sku struct {
	gorm.Model
	Sku_type_id uint64
	Item_id     uint64
	Code        string `gorm:"size:100;not null"`
}
