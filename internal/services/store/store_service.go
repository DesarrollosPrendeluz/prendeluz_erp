package services

import (
	"io"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
)

type StoreService interface {
	UpdateStoreStock(order []models.OrderItem) error
	GetStoreStock(storeName string, page int, pageSize int, searchParam string) []dtos.ItemStockInfo
	GetParent(child uint64) (models.Item, error)
	UploadStocks(file io.Reader, filename string) ([]StockUpdateError, error)
}
