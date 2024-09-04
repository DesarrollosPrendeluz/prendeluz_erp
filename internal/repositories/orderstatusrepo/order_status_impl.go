package orderstatusrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type OrderStatusImpl struct {
	*repositories.GORMRepository[models.OrderStatus]
}

func NewOrderStatussitory(db *gorm.DB) *OrderStatusImpl {
	return &OrderStatusImpl{repositories.NewGORMRepository(db, models.OrderStatus{})}
}
