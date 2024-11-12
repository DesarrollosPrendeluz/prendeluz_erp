package models

import "time"

type StoreLocation struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	StoreID   uint64
	Code      string
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (StoreLocation) TableName() string {
	return "store_locations"
}
