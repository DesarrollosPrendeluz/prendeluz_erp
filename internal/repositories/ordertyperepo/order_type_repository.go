package ordertyperepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type ItemRepo interface {
	repositories.Repository[models.OrderType]
}
