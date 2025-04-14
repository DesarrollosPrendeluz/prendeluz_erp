package services

import (
	"errors"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/itemlocationrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/storelocationrepo"
	"prendeluz/erp/internal/repositories/storestockrepo"
	"strconv"
)

type ItemStockLocationServiceImpl struct {
	itemlocationrepo  itemlocationrepo.ItemLocationImpl
	storelocationrepo storelocationrepo.StoreLocationImpl
	storestockrepo    storestockrepo.StoreStockRepoImpl
	orderitemrepo     orderitemrepo.OrderItemRepoImpl
	orderrepo         orderrepo.OrderRepoImpl
	itemrepo          itemsrepo.ItemRepoImpl
}

func NewItemStockLocationService() *ItemStockLocationServiceImpl {

	itemlocationrepo := *itemlocationrepo.NewInItemLocationRepository(db.DB)
	storelocationrepo := *storelocationrepo.NewStoreLocationRepository(db.DB)
	storestockrepo := *storestockrepo.NewStoreStockRepository(db.DB)
	orderitemrepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	orderrepo := *orderrepo.NewOrderRepository(db.DB)
	itemrepo := *itemsrepo.NewItemRepository(db.DB)

	return &ItemStockLocationServiceImpl{
		itemlocationrepo:  itemlocationrepo,
		storelocationrepo: storelocationrepo,
		storestockrepo:    storestockrepo,
		orderitemrepo:     orderitemrepo,
		orderrepo:         orderrepo,
		itemrepo:          itemrepo,
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
func (s *ItemStockLocationServiceImpl) GetItemStockLocation(main_sku string, store_id int, storeLocation int, page int, pageSize int) ([]models.ItemLocation, int64, error) {
	var err error
	var data []models.ItemLocation
	var datum *models.ItemLocation
	var recount int64

	if main_sku != "" && store_id != 0 {
		data, err = s.itemlocationrepo.FindByItemsAndStore(main_sku, uint64(store_id), pageSize, page)

	} else if main_sku != "" {
		data, err = s.itemlocationrepo.FindByItem(main_sku, pageSize, page)
	} else {
		if storeLocation != 0 {
			datum, err = s.itemlocationrepo.FindByID(uint64(storeLocation))
			if datum != nil {
				data = append(data, *datum)
			}
			recount = 1
		} else {
			data, err = s.itemlocationrepo.FindAll(pageSize, page)
			recount, _ = s.itemlocationrepo.CountAll()

		}

	}
	return data, recount, err
}

func (s *ItemStockLocationServiceImpl) PostItemStockLocation(requestBody dtos.ItemStockLocationCreateReq) []uint64 {
	var idArray []uint64

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {

		model := models.ItemLocation{
			ItemMainSku:     dataItem.ItemMainSku,
			StoreLocationID: dataItem.StoreLocationID,
			Stock:           dataItem.Stock,
		}
		s.itemlocationrepo.Create(&model)
		idArray = append(idArray, model.ID)
	}
	return idArray

}

func (s *ItemStockLocationServiceImpl) PatchItemStockLocation(requestBody dtos.ItemStockLocationUpdateReq) []error {
	var errorList []error

	for _, requestObject := range requestBody.Data {
		model, err := s.itemlocationrepo.FindByID(requestObject.Id)
		if err != nil {
			errorList = append(errorList, err)
			return errorList
		}
		if requestObject.ItemMainSku != nil {
			model.ItemMainSku = *requestObject.ItemMainSku
		}
		if requestObject.StoreLocationID != nil {
			model.StoreLocationID = *requestObject.StoreLocationID
		}
		if requestObject.Stock != nil {
			model.Stock = *requestObject.Stock
		}
		error := s.itemlocationrepo.Update(model)
		if error != nil {
			errorList = append(errorList, error)
		}

	}
	return errorList

}

func (s *ItemStockLocationServiceImpl) StockChanges(requestBody dtos.ItemStockLocationStockChangeRequest) []error {
	var errorList []error

	for _, requestObject := range requestBody.Data {
		model, err := s.itemlocationrepo.FindByID(requestObject.Id)
		loc, err1 := s.storelocationrepo.FindByID(model.StoreLocationID)
		stock, err2 := s.storestockrepo.FindByItemAndStore(model.ItemMainSku, strconv.FormatUint(loc.StoreID, 10))
		if err != nil || err1 != nil || err2 != nil {
			errorList = append(errorList, err)
			return errorList
		}

		stock.Amount = ((stock.Amount - int64(model.Stock)) + int64(requestObject.Stock))
		model.Stock = requestObject.Stock

		if stock.ReservedAmount != 0 && stock.ReservedAmount >= stock.Amount {
			stock.ReservedAmount = stock.Amount
			// //Update Picking (NOT GONNA BE USED)
			// items, _ := s.itemrepo.FindByEan(stock.Item.EAN)
			// var itemsIds []uint64
			// for _, item := range items {
			// 	itemsIds = append(itemsIds, item.ID)
			// }
			// ordersPicking := s.orderitemrepo.FindOrderByIteminPicking(itemsIds)
			// //Only will update the picking if its not complete
			// if ordersPicking.RecivedAmount == ordersPicking.Amount {
			// 	fmt.Println("No picking to update")
			// } else if ordersPicking.RecivedAmount == 0 && ordersPicking.Amount >= stock.ReservedAmount {
			// 	s.orderitemrepo.UpdatePickingByItemIdAndOrder(ordersPicking.ItemID, ordersPicking.OrderID, int(stock.Amount))
			// }
		}
		if requestObject.Stock >= 0 {
			error := s.itemlocationrepo.Update(model)
			error2 := s.storestockrepo.Update(&stock)
			if error != nil && error2 != nil {
				errorList = append(errorList, error)
				errorList = append(errorList, error2)
			}

		} else {

			errorList = append(errorList, errors.New("stock can't be negative"))
		}

	}
	return errorList

}

func (s *ItemStockLocationServiceImpl) StockMovements(requestBody dtos.ItemStockLocationStockMovementRequest) []error {
	var errorList []error

	stockMov := func(sku string, location uint64, stockVariant int64) error {
		model, err := s.itemlocationrepo.FindByItemsAndLocation(sku, location)
		loc, err1 := s.storelocationrepo.FindByID(model.StoreLocationID)
		stock, err2 := s.storestockrepo.FindByItemAndStore(model.ItemMainSku, strconv.FormatUint(loc.StoreID, 10))
		if err != nil || err1 != nil || err2 != nil {
			return errors.New("ha habido un error en la validacion del callback de stock movements")
		}

		stock.Amount = (stock.Amount + stockVariant)
		model.Stock = model.Stock + int(stockVariant)
		if model.Stock < 0 {
			return errors.New("stock can't be negative")

		}
		error := s.itemlocationrepo.Update(&model)
		error2 := s.storestockrepo.Update(&stock)
		if error != nil || error2 != nil {
			errorList = append(errorList, error)
			errorList = append(errorList, error2)
		}
		return nil
	}

	// Acceder a los valores del cuerpo

	for _, requestObject := range requestBody.Data {
		if requestObject.StoreLocationID1 != requestObject.StoreLocationID2 {
			errMov := stockMov(requestObject.MainSku, requestObject.StoreLocationID1, -int64(requestObject.Stock))
			if errMov == nil {
				errMov2 := stockMov(requestObject.MainSku, requestObject.StoreLocationID2, int64(requestObject.Stock))
				if errMov2 != nil {
					errorList = append(errorList, errMov2)
				}
			} else {
				errorList = append(errorList, errMov)
			}

		} else {
			errorList = append(errorList, errors.New("the locations are the same"))
		}

	}
	return errorList

}
