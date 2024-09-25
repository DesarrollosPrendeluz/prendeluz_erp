package stockdeficitrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type StockDeficitImpl struct {
	*repositories.GORMRepository[models.StockDeficit]
}

func NewStockDeficitRepository(db *gorm.DB) *StockDeficitImpl {
	return &StockDeficitImpl{repositories.NewGORMRepository(db, models.StockDeficit{})}
}

func (repo *StockDeficitImpl) GetallByStore(storeId int, pageSize int, offset int) ([]models.StockDeficit, error) {
	var models []models.StockDeficit
	repo.DB.Preload("Item").Where("store_id = ?", storeId).Limit(pageSize).Offset(offset).Find(&models)
	return models, nil
}

func (repo *StockDeficitImpl) CountConditional(storeId int) (int64, error) {
	var count int64
	err := repo.DB.Table("stock_deficits").Count(&count).Where("store_id = ?", storeId).Error
	return count, err
}
