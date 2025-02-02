package services

import (
	"errors"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/itemlocationrepo"
)

type ItemStockLocationServiceImpl struct {
	itemlocationrepo itemlocationrepo.ItemLocationImpl
}

func NewItemStockLocationService() *ItemStockLocationServiceImpl {

	itemlocationrepo := *itemlocationrepo.NewInItemLocationRepository(db.DB)

	return &ItemStockLocationServiceImpl{
		itemlocationrepo: itemlocationrepo,
	}
}

func (s *ItemStockLocationServiceImpl) DropItemLocation(locationId uint64) error {
	model, error := s.itemlocationrepo.FindByID(locationId)
	if error != nil {
		return error
	}
	if model.Stock == 0 {
		s.itemlocationrepo.Delete(model.ID)

	} else {
		return errors.New("el stock de la ubicación no es 0")
	}
	return nil

}
