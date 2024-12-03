package orderitemrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type OrderItemRepoImpl struct {
	*repositories.GORMRepository[models.OrderItem]
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepoImpl {
	return &OrderItemRepoImpl{repositories.NewGORMRepository(db, models.OrderItem{})}
}

// Se retornan las lineas de un pedido por el id del pedido
func (repo *OrderItemRepoImpl) FindByOrder(idOrder uint64) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	results := repo.DB.Preload("AssignedRel").Preload("AssignedRel.UserRel").Where("order_id = ?", idOrder).Find(&orderItems)

	return orderItems, results.Error
}

// Se retornan las lineas de pedidos las en las cuales el id de los items coincidan
func (repo *OrderItemRepoImpl) FindByItem(idPedido uint64) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	results := repo.DB.Where("item_id = ?", idPedido).Find(&orderItems)

	return orderItems, results.Error
}

// Se retornan las lineas de pedidos las en las cuales el id de los items coincidan
func (repo *OrderItemRepoImpl) FindByItemAndOrder(itemId uint64, orderId uint64) (models.OrderItem, error) {
	var orderItems models.OrderItem
	results := repo.DB.Where("item_id = ? and order_id = ?", itemId, orderId).First(&orderItems)

	return orderItems, results.Error
}
