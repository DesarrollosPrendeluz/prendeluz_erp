package models

import "time"

type ErpUpdateOrderLineHistory struct {
	ID                 uint   `gorm:"primaryKey"`
	UpdateErpTypeID    uint   `gorm:"index;not null"`
	UpdateGroupCode    string `gorm:"column:update_group_code"`
	OrderLineID        uint   `gorm:"index;not null"`
	OrderID            uint   `gorm:"index;not null"`
	ItemID             uint   `gorm:"index;not null"`
	UserID             uint   `gorm:"index;not null"`
	StoreID            uint   `gorm:"index;not null"`
	Quantity           int    `gorm:"not null;column:quantity"`
	NewQuantity        int    `gorm:"not null;column:new_quantity"`
	RecivedQuantity    int    `gorm:"default:0;not null;column:recived_quantity"`
	NewRecivedQuantity int    `gorm:"default:0;not null;column:new_recived_quantity"`
	CreatedAt          time.Time
	UpdatedAt          time.Time

	// Relaciones
	UpdateErpType UpdateErpType `gorm:"foreignKey:UpdateErpTypeID"`
	OrderLine     OrderItem     `gorm:"foreignKey:OrderLineID"`
	Order         Order         `gorm:"foreignKey:OrderID"`
	Item          Item          `gorm:"foreignKey:ItemID"`
	User          User          `gorm:"foreignKey:UserID"`
	Store         Store         `gorm:"foreignKey:StoreID"`
}

func (ErpUpdateOrderLineHistory) TableName() string {
	return "erp_update_order_line_histories"
}
