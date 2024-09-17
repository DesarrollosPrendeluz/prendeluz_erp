package orderstatusrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrderStatus struct {
	repositories.Repository[models.OrderStatus]
}
