package stockdeficitrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type StockDeficitRepo interface {
	repositories.Repository[models.StockDeficit]
	GetallByStore(storeId int, pageSize int, offset int) ([]models.StockDeficit, error)
	CountConditional(storeId int) (int64, error)
	GetByRegsitersByFatherSkuIn(filter []string, store int, page int, pageSize int) ([]models.StockDeficit, error)
	GetByFatherAndStore(fatherSku string, store int64) (models.StockDeficit, error)
	StockDeficitByFatherOrder(father_id uint64) (*[]StockDeficitResult, error)
	CallStockDefProc()
	CallPendingStockProc()
}
