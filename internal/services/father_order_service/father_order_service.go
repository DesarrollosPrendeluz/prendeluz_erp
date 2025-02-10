package services

import "prendeluz/erp/internal/dtos"

type FatherOrderService interface {
	FindLinesByFatherOrderCode(pageSize int, offset int, fatherOrderCode string, ean string, supplier_sku string, storeId int, searchByEan string, searchByLoc string, locFilter string) (dtos.FatherOrderOrdersAndLines, int64, error)
	DownloadOrdersExcelToAmazon(fatherCode string) string
	DownloadExcelAmazon(fatherID uint64) string
	CreateOrder(requestBody dtos.OrderWithLinesRequest) bool
}
