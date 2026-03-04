package asinrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type AsinRepoImpl struct {
	*repositories.GORMRepository[models.Asin]
}

func NewAsinRepository(db *gorm.DB) *AsinRepoImpl {
	return &AsinRepoImpl{repositories.NewGORMRepository(db, models.Asin{})}
}

func (repo *AsinRepoImpl) FindByItemId(id uint64) (models.Asin, error) {
	var asin models.Asin

	result := repo.DB.Where("item_id = ?", id).First(&asin)

	return asin, result.Error

}

func (repo *AsinRepoImpl) FindByAsins(asin []string) ([]models.Asin, error) {
	var asins []models.Asin

	result := repo.DB.Where("code IN ?", asin).Find(&asins)

	return asins, result.Error
}
