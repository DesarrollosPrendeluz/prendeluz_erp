package services

import (
	dtos "prendeluz/erp/internal/dtos/api"
	"prendeluz/erp/internal/models"
)

type StockService interface {
	ReturnDownloadStockExcel(store_id int) string
	FindStockItems(itemId uint64) []models.Item
	ReturnStockByAsins(asins []string) []dtos.StockItem
}
