package services

import (
	"io"
	"prendeluz/erp/internal/dtos"
)

type OrderService interface {
	UploadOrderExcel(file io.Reader, filename string) error
	GetOrders(page int, pageSize int, startDate string, endDate string, typeId int, statusId int, orderCode string) ([]dtos.ItemsPerOrder, int64, error)
	OrderComplete(orderCode string) error
}
