package storerepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type StoreRepo interface {
	repositories.Repository[models.Store]
	FindByName(name string) models.Store
}
