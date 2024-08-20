package repositories

import (
	"gorm.io/gorm"
	"prendeluz/erp/internal/models"
)

type OrderRepo struct {
	GORMRepository[models.Order]
}

func NewOrderRepository(db *gorm.DB) *OrderRepo {
	return &OrderRepo{*NewGORMRepository(db, models.Order{})}

}
