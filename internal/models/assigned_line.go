package models

import "time"

type AssignedLine struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	UserID      uint64 `gorm:"column:user_id;not null"`
	OrderLineID uint64 `gorm:"column:order_line_id;not null"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (AssignedLine) TableName() string {
	return "assigned_lines"
}
