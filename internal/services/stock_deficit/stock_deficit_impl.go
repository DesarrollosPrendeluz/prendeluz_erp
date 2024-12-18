package services

import (
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/stockdeficitrepo"
)

type StockDeficitServiceImpl struct {
	stockDeficitRepo stockdeficitrepo.StockDeficitImpl
	itemsRepo        itemsrepo.ItemRepoImpl
	orderErrorRepo   repositories.GORMRepository[models.ErrorOrder]
}

func NewStockDeficitService() *StockDeficitServiceImpl {
	stockDeficitRepo := *stockdeficitrepo.NewStockDeficitRepository(db.DB)
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)

	return &StockDeficitServiceImpl{
		stockDeficitRepo: stockDeficitRepo,
		orderErrorRepo:   errorOrderRepo,
		itemsRepo:        itemsRepo}
}

func (s *StockDeficitServiceImpl) SearchBySkuAndEan(filter string, store int, page int, pageSize int) ([]models.StockDeficit, error) {

	var modelsData []models.StockDeficit

	subQuery := s.orderErrorRepo.DB.
		Model(&models.StockDeficit{}).
		Select("stock_deficits.id").
		Joins("JOIN items ON items.main_sku = stock_deficits.parent_main_sku").
		Where("store_id = ?", store).
		Where("items.main_sku LIKE ?", "%"+filter+"%").
		Or("items.ean LIKE ?", "%"+filter+"%")

	err := s.orderErrorRepo.DB.Debug().
		Preload("Item.SupplierItems.Supplier").
		Where("id IN (?)", subQuery).
		Limit(pageSize).
		Offset(page).
		Find(&modelsData).Error

	if err != nil {
		return nil, err
	}
	return modelsData, nil

}
