package itemlocationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type StoreLocationRepo interface {
	repositories.Repository[models.StoreLocation]
	FindByItem(mainSku string, pageSize int, offset int) ([]models.ItemLocation, error)
	FindByItemsAndStore(mainSku string, store uint64, pageSize int, offset int) ([]models.ItemLocation, error)
	FindByItemsAndLocation(mainSku string, location uint64) (models.ItemLocation, error)
	FindByItemAndLocation(mainSku string, locationId uint64) (models.ItemLocation, error)
	DeleteZeroStock() error
}
