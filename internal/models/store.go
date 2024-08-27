package models

type Store struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:255,not null"`
}

func (Store) TableName() string {
	return "store"
}
