package storelocationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type StoreLocationRepo interface {
	repositories.Repository[models.StoreLocation]
	FindStoreLocationByCode(code string) (models.StoreLocation, error)
}
