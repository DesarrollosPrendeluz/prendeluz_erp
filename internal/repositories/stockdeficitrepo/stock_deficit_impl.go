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
