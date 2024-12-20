package models

import (
	"time"
)

// OrderLineBox represents a box in an order line
type OrderLineBox struct {
	ID          int       `json:"id" gorm:"column:id"`
	OrderLineID int       `json:"order_line_id" gorm:"column:order_line_id"`
	BoxID       string    `json:"box_id" gorm:"column:box_id"`
	Quantity    int       `json:"quantity" gorm:"column:quantity"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`

	Boxes      *[]Box       `gorm:"foreignKey:ID;references:BoxID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrderLines *[]OrderItem `gorm:"foreignKey:ID;references:OrderLineID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (OrderLineBox) TableName() string {
	return "order_lines_boxes"
}
