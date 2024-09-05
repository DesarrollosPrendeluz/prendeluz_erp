package orderrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrdeRepo interface {
	repositories.Repository[models.Order]
	FindOrderByDate(startDate string, endDate string) (models.Order, error)
	FindByOrderCode(orderCode string) (models.Order, error)
	UpdateStatus(newStatus string, orderID uint64) error
}
