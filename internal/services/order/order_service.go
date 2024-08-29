package services

import (
	"io"
	"prendeluz/erp/internal/dtos"
)

type OrderService interface {
	UploadOrderExcel(file io.Reader, filename string) error
	GetOrders(page int, pageSize int) ([]dtos.ItemsPerOrder, error)
	OrderComplete(orderCode string) error
}
