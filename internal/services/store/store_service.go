package services

import (
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
)

type StoreService interface {
	UpdateStoreStock(order []models.OrderItem) error
	GetStoreStock() []dtos.StoreStockItems
}
