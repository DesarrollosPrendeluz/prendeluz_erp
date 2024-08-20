package repositories

import (
	"prendeluz/erp/internal/models"

	"gorm.io/gorm"
)

type OrderItemRepo struct {
	GORMRepository[models.OrderItem]
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepo {
	return &OrderItemRepo{*NewGORMRepository(db, models.OrderItem{})}
}

func (repo *OrderItemRepo) FindByOrder(idOrder uint64) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	results := repo.db.Where("id_pedido = ?", idOrder).Find(&orderItems)

	return orderItems, results.Error
}
