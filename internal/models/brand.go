package models

import "time"

type Brand struct {
	ID        uint64 `gorm:"primary_key;not null"`
	Name      string `gorm:"column:name;not null"`
	Email     string `gorm:"column:email;not null"`
	Address   string `gorm:"column:address;not null"`
	Company   string `gorm:"column:company;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Brand) TableName() string {
	return "brands"
}
