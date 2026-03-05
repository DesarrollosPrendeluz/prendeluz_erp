package services

import (
	"io"
	"prendeluz/erp/internal/dtos"
	apiDtos "prendeluz/erp/internal/dtos/api"
)

type OrderService interface {
	UploadOrderExcel(file io.Reader, filename string) error
	GetOrders(page int, pageSize int, startDate string, endDate string, typeId int, statusId int, orderCode string) ([]dtos.ItemsPerOrder, int64, error)
	OrderComplete(orderCode string) error
	CreateOrderViaAPI(order apiDtos.ApiOrderCreate) error
}
