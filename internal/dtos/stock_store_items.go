package dtos

import "prendeluz/erp/internal/models"

type ItemStockInfo struct {
	Itemname string
	SKU      string
	Childs   []models.Item
	Amount   int64
}
