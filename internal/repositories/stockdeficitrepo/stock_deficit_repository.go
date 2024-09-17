package stockdeficitrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type OrderTypeImpl struct {
	*repositories.GORMRepository[models.StockDeficit]
}

func NewStockDeficitRepository(db *gorm.DB) *OrderTypeImpl {
	return &OrderTypeImpl{repositories.NewGORMRepository(db, models.StockDeficit{})}
}
