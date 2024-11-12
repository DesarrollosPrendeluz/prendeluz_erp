package models

import "time"

type ItemLocation struct {
	ID              uint64 `gorm:"primaryKey;autoIncrement"`
	ItemMainSku     string `gorm:"size:255;not null;column:item_main_sku"`
	StoreLocationID uint64 `gorm:"not null;column:store_location_id"`
	Stock           int
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	StoreLocations  *StoreLocation `gorm:"foreignKey:ID;references:StoreLocationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ItemLocation) TableName() string {
	return "item_stock_locations"
}
