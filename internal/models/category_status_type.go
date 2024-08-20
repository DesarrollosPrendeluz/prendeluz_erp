package models

import "time"

type CategoryStatusType struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:255;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
