package orderrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrdeRepo interface {
	repositories.Repository[models.Order]

	FindByOrderCode(orderCode string) (models.Order, error)
	UpdateStatus(newStatus string, orderID uint64) error
}
