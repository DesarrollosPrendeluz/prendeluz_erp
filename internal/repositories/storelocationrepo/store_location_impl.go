package storelocationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type StoreLocationImpl struct {
	*repositories.GORMRepository[models.StoreLocation]
}

func NewStoreLocationRepository(db *gorm.DB) *StoreLocationImpl {
	return &StoreLocationImpl{repositories.NewGORMRepository(db, models.StoreLocation{})}
}

func (repo *StoreLocationImpl) FindStoreLocationByCode(code string) (models.StoreLocation, error) {
	var StoreLoc models.StoreLocation

	result := repo.DB.
		Where("code = ? ", code).
		First(&StoreLoc)
	return StoreLoc, result.Error
}
