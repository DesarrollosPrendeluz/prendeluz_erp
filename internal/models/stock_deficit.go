package models

type StockDeficit struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement"`
	OrderID uint64 `gorm:"column:id_order;primaryKey;not null"`
	ItemId  string `gorm:"column:id_item;primaryKey;not null"`
	Amount  int64

	Item  Item  `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Store Store `gorm:"foreignKey:StoreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
