package repositories

import (
	"prendeluz/erp/internal/models"

	"gorm.io/gorm"
)

type ErrorOrderRepo struct {
	GORMRepository[models.ErrorOrder]
}

func NewErrorOrderRepository(db *gorm.DB) *ErrorOrderRepo {
	return &ErrorOrderRepo{*NewGORMRepository(db, models.ErrorOrder{})}
}
