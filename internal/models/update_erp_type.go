package models

import "time"

type UpdateErpType struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:255,not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (UpdateErpType) TableName() string {
	return "update_erp_types"
}
