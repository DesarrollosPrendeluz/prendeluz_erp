package models

import "time"

type AssignedLine struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	UserID      uint64 `gorm:"column:user_id;not null"`
	OrderLineID uint64 `gorm:"column:order_line_id;not null"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time

	UserRel User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (AssignedLine) TableName() string {
	return "assigned_lines"
}
