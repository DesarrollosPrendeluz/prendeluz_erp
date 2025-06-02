package services

import "prendeluz/erp/internal/models"

type StockService interface {
	ReturnDownloadStockExcel(store_id int) string
	FindStockItems(itemId uint64) []models.Item
}
