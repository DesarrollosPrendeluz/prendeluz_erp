package dtos

type ItemStockInfo struct {
	itemname string
	sku      string
	cantidad int64
}
type StoreStockItems struct {
	name  string
	items []ItemStockInfo
}
