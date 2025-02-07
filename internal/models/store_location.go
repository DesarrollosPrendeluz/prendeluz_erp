package models

import "time"

type StoreLocation struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	StoreID   uint64
	Code      string
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time

	Store Store `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (StoreLocation) TableName() string {
	return "store_locations"
}
