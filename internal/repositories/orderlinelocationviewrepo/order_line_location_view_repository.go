package orderlinelocationviewrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrderLineLocationViewRepo interface {
	repositories.Repository[models.OrderLineLocationView]
	FindByFatherAndStoreWithOrder(father_id uint64, idStore uint64, orderByLocation string, orderByEan string) ([]uint64, string, error)
}
