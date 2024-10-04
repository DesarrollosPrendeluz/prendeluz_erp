package models

type SupplierItem struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	SupplierID  uint64
	BrandID     uint64
	ItemID      uint64
	Order       int
	SupplierSku string `gorm:"size:255,not null"`

	Brand *Brand `gorm:"foreignKey:BrandId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (SupplierItem) TableName() string {
	return "supplier_items"
}
