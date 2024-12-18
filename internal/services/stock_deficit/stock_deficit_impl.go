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

func (s *StockDeficitServiceImpl) SearchBySkuAndEan(filter string, store int, page int, pageSize int) ([]models.StockDeficit, []error) {

	var fatherSkus []string
	var errArray []error
	//subQuery := s.stockDeficitRepo.
	items, err1 := s.itemsRepo.FindByFathersMainSkuOrEan(filter)
	for _, item := range items {
		fatherSkus = append(fatherSkus, item.MainSKU)
	}

	stockDef, err2 := s.stockDeficitRepo.GetByRegsitersByFatherSkuIn(fatherSkus, store, page, pageSize)

	if err1 != nil || err2 != nil {
		errArray = append(errArray, err1)
		errArray = append(errArray, err2)
		return nil, errArray
	}
	return stockDef, errArray

}
