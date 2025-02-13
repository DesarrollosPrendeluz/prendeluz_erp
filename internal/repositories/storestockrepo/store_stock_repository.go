package storestockrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type StoreStockRepo interface {
	repositories.Repository[models.StoreStock]
	FindByItem(parent_sku string) (models.StoreStock, error)
	FindByStore(idStore uint64, pageSize int, offset int) ([]models.StoreStock, error)
	FindByStoreAndSearchParams(idStore uint64, searchParam string, pageSize int, offset int) ([]models.StoreStock, error)
	FindByStoreWithLocations(idStore uint64) (StoreAndlocations, error)
}
