package models

import "time"

type StockDeficit struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	StoreID    uint64 `gorm:"column:store_id;primaryKey;not null"`
	SKU_Parent string `gorm:"column:parent_main_sku;primaryKey;not null"`
	Amount     int64  `gorm:"column:quantity"`
	Item       Item   `gorm:"foreignKey:MainSKU;references:SKU_Parent;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Store      Store  `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func (StockDeficit) TableName() string {
	return "stock_deficits"
}
