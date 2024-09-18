package models

import "time"

type OrderItem struct {
	ID            uint64 `gorm:"autoIncrement;primaryKey"`
	OrderID       uint64 `gorm:"column:order_id;not null"`
	ItemID        uint64 `gorm:"column:item_id;not null"`
	Amount        int64  `gorm:"column:quantity;not null"`
	RecivedAmount int64  `gorm:"column:recived_quantity;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Order Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Item  Item  `gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (OrderItem) TableName() string {
	return "order_lines"
}
