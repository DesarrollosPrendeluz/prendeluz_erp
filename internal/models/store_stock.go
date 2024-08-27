package models

type StoreStock struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	StoreID    uint64 `gorm:"column:id_store;primaryKey;not null"`
	SKU_Parent string `gorm:"column:parent_sku;primaryKey;not null"`
	Amount     int64

	Item  Item  `gorm:"foreignKey:SKU_Parent;references:main_sku;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Store Store `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (StoreStock) TableName() string {
	return "store_stock"
}
