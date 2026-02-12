package services

import (
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
)

type ItemStockLocationService interface {
	DropItemLocation(locationId uint64) error
	GetItemStockLocation(main_sku string, store_id int, storeLocation int, page int, pageSize int) ([]models.ItemLocation, int64, error)
	PostItemStockLocation(requestBody dtos.ItemStockLocationCreateReq) []uint64
	PatchItemStockLocation(requestBody dtos.ItemStockLocationUpdateReq) []error
	StockChanges(requestBody dtos.ItemStockLocationStockChangeRequest) []error
	StockMovements(requestBody dtos.ItemStockLocationStockMovementRequest) []error
	DeleteZeroStock() error
}
