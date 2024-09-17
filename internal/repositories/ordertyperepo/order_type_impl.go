package ordertyperepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type OrderTypeImpl struct {
	*repositories.GORMRepository[models.OrderType]
}

func NewOrderTypeRepository(db *gorm.DB) *OrderTypeImpl {
	return &OrderTypeImpl{repositories.NewGORMRepository(db, models.OrderType{})}
}
