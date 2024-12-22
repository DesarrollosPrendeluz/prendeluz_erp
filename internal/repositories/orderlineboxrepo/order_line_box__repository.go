package storelocationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrderLineBoxRepo interface {
	repositories.Repository[models.OrderLineBox]
}
