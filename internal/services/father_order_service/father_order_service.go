package services

import "prendeluz/erp/internal/dtos"

type FatherOrderService interface {
	FindLinesByFatherOrderCode(pageSize int, offset int, fatherOrderCode string, ean string, supplier_sku string, storeId int) (dtos.FatherOrderOrdersAndLines, int64, error)
	DownloadOrdersExcelToAmazon(fatherCode string) string
}
