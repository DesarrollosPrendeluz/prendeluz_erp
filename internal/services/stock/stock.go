package services

type StockService interface {
	ReturnDownloadStockExcel(store_id int) string
}
