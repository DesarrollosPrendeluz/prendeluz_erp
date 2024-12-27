package orderitemrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrderItemRepo interface {
	repositories.Repository[models.OrderItem]
	FindByOrder(idOrder uint64) ([]models.OrderItem, error)
	FindByItem(idPedido uint64) ([]models.OrderItem, error)
	FindByOrderAndItem(orderIds []uint64, storeId int, itemIds []uint64, offset int, pageSize int) ([]models.OrderItem, int64)
}
