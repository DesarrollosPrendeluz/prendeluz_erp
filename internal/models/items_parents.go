package models

import "time"

type ItemsParents struct {
	ID           uint64
	ChildItemID  uint64
	ParentItemID uint64
	CreatedAt    *time.Time
	UpdatedAt    *time.Time

	Child  *Item `gorm:"foreignKey:ChildItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Parent *Item `gorm:"foreignKey:ParentItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ItemsParents) TableName() string {
	return "item_parents"
}
