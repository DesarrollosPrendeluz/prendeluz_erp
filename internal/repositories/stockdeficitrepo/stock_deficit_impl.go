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
	repo.DB.Preload("Item.SupplierItems.Supplier").Where("store_id = ?", storeId).Limit(pageSize).Offset(offset).Find(&models)
	return models, nil
}
func (repo *StockDeficitImpl) GetallByStoreAndSupplier(storeId int, supplier int, pageSize int, offset int) ([]models.StockDeficit, error) {
	var modelsData []models.StockDeficit

	subQuery := repo.DB.
		Model(&models.StockDeficit{}).
		Select("stock_deficits.id").
		Joins("JOIN items ON items.main_sku = stock_deficits.parent_main_sku").
		Joins("JOIN supplier_items ON supplier_items.item_id = items.id").
		//Where("supplier_items.item_id = ?", 2).
		Where("store_id = ?", storeId).
		Where("supplier_items.supplier_id = ?", supplier)

	err := repo.DB.
		//("Item.SupplierItems", "item_id = ?", 2).
		Preload("Item.SupplierItems.Supplier").
		Where("id IN (?)", subQuery).
		Limit(pageSize).
		Offset(offset).
		Find(&modelsData).Error

	if err != nil {
		return nil, err
	}
	return modelsData, nil
}

func (repo *StockDeficitImpl) CountConditional(storeId int) (int64, error) {
	var count int64
	err := repo.DB.Table("stock_deficits").Count(&count).Where("store_id = ?", storeId).Error
	return count, err
}
