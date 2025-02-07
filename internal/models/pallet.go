package models

import (
	"time"
)

type Pallet struct {
	ID      uint64 `json:"id" gorm:"primary_key;column:id"`
	OrderID uint64 `json:"order_id" gorm:"column:order_id"`
	Number  int    `json:"number" gorm:"column:number"`
	Label   string `json:"label" gorm:"column:label"`
	IsClose int    `gorm:"column:is_close;not null;default:0"`

	Order *Order `gorm:"foreignKey:ID;references:OrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Boxes *[]Box `gorm:"foreignKey:PalletID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Pallet) TableName() string {
	return "pallets"
}
