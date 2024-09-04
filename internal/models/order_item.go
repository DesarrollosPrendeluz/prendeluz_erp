package models

type OrderItem struct {
	ID      uint64 `gorm:"autoIncrement"`
	OrderID uint64 `gorm:"column:order_id;primaryKey;not null"`
	ItemID  uint64 `gorm:"column:item_id;primaryKey;not null"`
	Amount  int64  `gorm:"column:quantity;not null"`

	Order Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Item  Item  `gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (OrderItem) TableName() string {
	return "order_lines"
}
