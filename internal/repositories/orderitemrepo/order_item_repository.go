package orderitemrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrderItemRepo interface {
	repositories.Repository[models.OrderItem]
}
