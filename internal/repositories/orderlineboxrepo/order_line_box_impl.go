package orderlineboxrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type OrderLineBoxImpl struct {
	*repositories.GORMRepository[models.OrderLineBox]
}

func NewOrderLineBoxRepository(db *gorm.DB) *OrderLineBoxImpl {
	return &OrderLineBoxImpl{repositories.NewGORMRepository(db, models.OrderLineBox{})}
}
