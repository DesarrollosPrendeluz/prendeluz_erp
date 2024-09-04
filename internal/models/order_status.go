package models

type OrderStatus struct {
	ID   uint64 `gorm:"autoIncrement"`
	Name string `gorm:"size:255;not null"`
}

func (OrderStatus) TableName() string {
	return "order_statuses"
}
