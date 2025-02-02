package services

type ItemStockLocationService interface {
	DropItemLocation(locationId uint64) error
}
