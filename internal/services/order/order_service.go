package services

import (
	"io"
	"prendeluz/erp/internal/dtos"
)

type OrderService interface {
	UploadOrderExcel(file io.Reader, filename string) error
	GetOrders() ([]dtos.ItemsPerOrder, error)
}
