package orderrepo

import (
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrdeRepo interface {
	repositories.Repository[models.Order]
	FindOrderFiltered(pageSize int, page int, startDate string, endDate string, typeId int, statusId int, orderCode string) ([]models.Order, int64, error)
	UpdateStatus(newStatus string, orderID uint64) error
	GetSupplierOrders(order_type *int) ([]dtos.SupplierOrders, error)
	GetSupplierOrdersByFatherSku(fatherOrderId int) ([]dtos.SupplierOrders, error)
}
