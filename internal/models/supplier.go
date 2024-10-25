package models

type Supplier struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"size:255,not null"`
	DeliveryTime string `gorm:"size:255,not null"`
}

func (Supplier) TableName() string {
	return "suppliers"
}
