package models

import (
	"time"
)

type Box struct {
	ID       uint64 `json:"id" gorm:"primary_key;column:id"`
	PalletID uint64 `json:"palletId" gorm:"column:pallet_id"`
	Number   int    `json:"number" gorm:"column:number"`
	Label    string `json:"label" gorm:"column:label"`
	Quantity int    `json:"quantity" gorm:"column:quantity"`
	IsClose  int    `gorm:"column:is_close;not null;default:0"`

	Pallet     *Pallet         `gorm:"foreignKey:ID;references:PalletID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BoxContent *[]OrderLineBox `gorm:"foreignKey:BoxID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Box) TableName() string {
	return "boxes"
}
