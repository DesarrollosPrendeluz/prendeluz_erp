package models

import "time"

type AccesTokens struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	UserId    uint64 `gorm:"primaryKey;autoIncrement"`
	Token     string `gorm:"size:255;not null;uniqueIndex"`
	Valid     bool   `gorm:"not null;default:false"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (AccesTokens) TableName() string {
	return "user_token_erps"
}
