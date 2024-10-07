package models

type SupplierItem struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	SupplierID  uint64
	BrandID     uint64
	ItemID      uint64
	Order       int
	SupplierSku string `gorm:"size:255,not null"`

	Item  *Item  `gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Brand *Brand `gorm:"foreignKey:ID;references:BrandID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (SupplierItem) TableName() string {
	return "supplier_items"
}
