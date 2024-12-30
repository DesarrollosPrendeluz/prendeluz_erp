package services

import (
	"prendeluz/erp/internal/models"
)

type StockDeficitService interface {
	SearchBySkuAndEan(filter string, store int, page int, pageSize int) ([]models.StockDeficit, []error)
	CalcStockDeficitByItem(child_item_id uint64, store_id int64)
}
