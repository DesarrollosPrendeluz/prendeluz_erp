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
	FindByItemsAndOrder(itemIds []uint64, orderId uint64) (models.OrderItem, error)
	FindByOrderAndStore(idOrder uint64, store_id int) ([]models.OrderItem, error)
	FindByLineIDWithOrder(lineId []uint64, order string, offset int, pageSize int) ([]models.OrderItem, int64)

}
