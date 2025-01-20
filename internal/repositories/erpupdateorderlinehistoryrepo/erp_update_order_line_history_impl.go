package erpupdateorderlinehistoryrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type ErpUpdateOrderLineHistory struct {
	*repositories.GORMRepository[models.ErpUpdateOrderLineHistory]
}

func NewErpUpdateOrderLineHistoryRepository(db *gorm.DB) *ErpUpdateOrderLineHistory {
	return &ErpUpdateOrderLineHistory{repositories.NewGORMRepository(db, models.ErpUpdateOrderLineHistory{})}
}
