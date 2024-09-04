package models

type Asin struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement"`
	BrandId uint64
	ItemID  uint64
	Code    string `gorm:"size:255,not null"`
	Ean     string `gorm:"size:255,not null"`
}

func (Asin) TableName() string {
	return "asins"
}
