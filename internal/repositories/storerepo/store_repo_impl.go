package storerepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type StoreRepoImpl struct {
	*repositories.GORMRepository[models.Store]
}

func NewStoreRepository(db *gorm.DB) *StoreRepoImpl {
	return &StoreRepoImpl{repositories.NewGORMRepository(db, models.Store{})}
}

func (repo *StoreRepoImpl) FindByName(name string) models.Store {
	var store models.Store

	repo.DB.Where("name LIKE ?", "%"+name+"%").First(&store)

	return store
}
