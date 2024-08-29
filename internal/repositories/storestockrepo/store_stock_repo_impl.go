package storestockrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type StoreStockRepoImpl struct {
	*repositories.GORMRepository[models.StoreStock]
}

func NewStoreStockRepository(db *gorm.DB) *StoreStockRepoImpl {
	return &StoreStockRepoImpl{repositories.NewGORMRepository(db, models.StoreStock{})}
}
func (repo *StoreStockRepoImpl) FindByItem(sku_parent string) (models.StoreStock, error) {
	var storeStocks models.StoreStock

	results := repo.DB.Where("parent_sku LIKE ?", "%"+sku_parent+"%").First(&storeStocks)

	return storeStocks, results.Error
}

func (repo *StoreStockRepoImpl) FindByStore(idStore uint64, pageSize int, offset int) ([]models.StoreStock, error) {
	var storeStocks []models.StoreStock

	results := repo.DB.Limit(pageSize).Offset(offset).Preload("Item").Where("id_store = ?", idStore).Find(&storeStocks)

	return storeStocks, results.Error

}
