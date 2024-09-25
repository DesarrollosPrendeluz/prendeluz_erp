package stockdeficitrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type StockDeficitRepo interface {
	repositories.Repository[models.StockDeficit]
	GetallByStore(storeId int, pageSize int, offset int) ([]models.StockDeficit, error)
	CountConditional(storeId int) (int64, error)
}
