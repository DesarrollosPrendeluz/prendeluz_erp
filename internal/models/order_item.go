package models

import "time"

type OrderItem struct {
	ID            uint64 `gorm:"autoIncrement;primaryKey"`
	OrderID       uint64 `gorm:"column:order_id;not null"`
	ItemID        uint64 `gorm:"column:item_id;not null"`
	Amount        int64  `gorm:"column:quantity;not null"`
	RecivedAmount int64  `gorm:"column:recived_quantity;not null"`
	StoreID       int64  `gorm:"column:store_id;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	ClientID uint64 `gorm:"-"`

	Order            Order            `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Item             Item             `gorm:"foreignKey:ID;references:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssignedRel      AssignedLine     `gorm:"foreignKey:OrderLineID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InOrderRelation  InOrderRelation  `gorm:"foreignKey:ID;references:OrderLineID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OutOrderRelation OutOrderRelation `gorm:"foreignKey:ID;references:OrderLineID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (OrderItem) TableName() string {
	return "order_lines"
}
