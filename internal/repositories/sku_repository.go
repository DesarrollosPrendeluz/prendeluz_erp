package repositories

import (
	"prendeluz/erp/internal/models"

	"gorm.io/gorm"
)

type SkuRepo struct {
	GORMRepository[models.Sku]
}

func NewSkuRepository(db *gorm.DB) *SkuRepo {
	return &SkuRepo{*NewGORMRepository(db, models.Sku{})}

}

func (repo *SkuRepo) FindByCode(code string) ([]models.Sku, error) {
	var skus []models.Sku
	results := repo.DB.Where("code LIKE ?", "%"+code+"%").Find(&skus)

	return skus, results.Error
}
