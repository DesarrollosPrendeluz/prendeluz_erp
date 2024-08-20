package models

import "time"

type Category struct {
	ID                 uint64 `gorm:"primaryKey;autoIncrement"`
	Name               string `gorm:"size:255;not null"`
	TypeOfCategoriesID uint64 `gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CategoryStatusID   uint64 `gorm:"not null"`

	CategoryStatusType CategoryStatusType `gorm:"foreignKey:CategoryStatusID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TypeOfCategories   TypeOfCategories   `gorm:"foreignKey:TypeOfCategoriesID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
