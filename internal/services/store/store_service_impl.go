package services

import (
	"io"
	"log"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/itemlocationrepo"
	"prendeluz/erp/internal/repositories/itemsparentsrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/storelocationrepo"
	"prendeluz/erp/internal/repositories/storerepo"
	"prendeluz/erp/internal/repositories/storestockrepo"
	"prendeluz/erp/internal/utils"
	"strconv"

	"gorm.io/gorm"
)

type StoreServiceImpl struct {
	orderRepo         *orderrepo.OrderRepoImpl
	orderItemsRepo    *orderitemrepo.OrderItemRepoImpl
	itemsRepo         *itemsrepo.ItemRepoImpl
	storeStockRepo    *storestockrepo.StoreStockRepoImpl
	itemsParentsRepo  *itemsparentsrepo.ItemsParentsRepoImpl
	storeRepo         *storerepo.StoreRepoImpl
	itemlocationrepo  *itemlocationrepo.ItemLocationImpl
	storelocationrepo *storelocationrepo.StoreLocationImpl
}

func NewStoreService() *StoreServiceImpl {
	orderRepo := orderrepo.NewOrderRepository(db.DB)
	orderItemRepo := orderitemrepo.NewOrderItemRepository(db.DB)
	itemsRepo := itemsrepo.NewItemRepository(db.DB)
	storeStockRepo := storestockrepo.NewStoreStockRepository(db.DB)
	itemsParentsRepo := itemsparentsrepo.NewItemParentRepository(db.DB)
	storeRepo := storerepo.NewStoreRepository(db.DB)
	itemlocationrepo := itemlocationrepo.NewInItemLocationRepository(db.DB)
	storelocationrepo := storelocationrepo.NewStoreLocationRepository(db.DB)

	return &StoreServiceImpl{orderRepo: orderRepo, orderItemsRepo: orderItemRepo, itemsRepo: itemsRepo, storeStockRepo: storeStockRepo, itemsParentsRepo: itemsParentsRepo, storeRepo: storeRepo, itemlocationrepo: itemlocationrepo, storelocationrepo: storelocationrepo}
}

// Obtiene un registro padre en base a uno de sus hijos
func (s *StoreServiceImpl) GetParent(child uint64) (models.Item, error) {
	itemsParent, _ := s.itemsParentsRepo.FindByChild(child)
	parent, err := s.itemsRepo.FindByID(itemsParent.ParentItemID)

	return *parent, err
}

// Actualiza el stock de un alamcen en base a una orden
func (s *StoreServiceImpl) UpdateStoreStock(orderCode string) error {
	itemsOrdered := make(map[string]int64)
	orders, _ := s.orderRepo.FindByOrderCode(orderCode)
	type StockDeficit struct {
		MainSku string
		Amount  int64
	}
	orderItems, _ := s.orderItemsRepo.FindByOrder(orders.ID)

	for _, order := range orderItems {
		item, _ := s.itemsRepo.FindByID(order.ItemID)
		parentSKU := item.MainSKU
		if item.ItemType != "father" {
			itemParent, err := s.GetParent(item.ID)
			parentSKU = itemParent.MainSKU
			if err != nil {
				log.SetPrefix("[ERROR]")
				log.Println("Parent not found for ", item.MainSKU, " SKU")
			}
		}

		itemsOrdered[parentSKU] += order.Amount
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		s.itemsRepo.SetDB(tx)
		s.orderItemsRepo.SetDB(tx)
		s.orderRepo.SetDB(tx)
		s.storeStockRepo.SetDB(tx)

		var updateStock []models.StoreStock

		for item, amount := range itemsOrdered {
			itemToUpdate, err := s.storeStockRepo.FindByItem(item)
			if err != nil {
				continue
			}
			itemToUpdate.Amount -= amount
			if itemToUpdate.Amount < 0 {
				deficit := -itemToUpdate.Amount
				sd := StockDeficit{MainSku: itemToUpdate.SKU_Parent, Amount: deficit}
				log.Println("Stock deficit: ", sd)
				itemToUpdate.Amount = 0
			}
			updateStock = append(updateStock, itemToUpdate)
		}
		s.storeStockRepo.UpdateAll(&updateStock)
		err := s.orderRepo.UpdateStatus(orderrepo.Order_Status["en_proceso"], orders.ID)
		log.Println(err)

		return nil
	})
}

// Obtiene los articulos hijos de un artículo padre
func getChilds(items []models.ItemsParents) []models.Item {
	var results []models.Item
	for _, child := range items {
		results = append(results, *child.Child)
	}
	return results
}

func (s *StoreServiceImpl) UploadStocks(file io.Reader, filename string) (string, string, error) {

	//fatherRepo := fatherorderrepo.NewFatherOrderRepository(db.DB)
	var stockErr []utils.StockUpdateError
	addError := func(errorData error, errArr *[]utils.StockUpdateError, sku string, loc string, err string) bool {
		if errorData != nil {
			errReturn := utils.StockUpdateError{
				FatherSku: sku,
				Loc:       loc,
				Error:     err,
			}
			*errArr = append(*errArr, errReturn)
			return false

		}
		return true

	}

	data, err := utils.ExcelToJsonUpdateStocks(file)

	if addError(err, &stockErr, "", "", "No se ha conseguido leer el archivo") {
		for _, datum := range data {
			fatherItem, err2 := returnFatherData(datum.Sku)

			if addError(err2, &stockErr, datum.Sku, datum.Loc, "No se ha encontrado el articulo padre") {
				loc, err3 := s.storelocationrepo.FindStoreLocationByCode(datum.Loc)

				if addError(err3, &stockErr, datum.Sku, datum.Loc, "Error la ubicación no ha sido encontrada") {
					itemLoc, _ := s.itemlocationrepo.FindByItemAndLocation(fatherItem.FatherSku, loc.ID)
					stock, _ := s.storeStockRepo.FindByItemAndStore(fatherItem.FatherSku, strconv.FormatUint(loc.StoreID, 10))
					stock.Amount = (stock.Amount - int64(itemLoc.Stock)) + datum.Quantity
					//if addError(err4, &stockErr, datum.Sku, datum.Loc, "Erroren encontrar el articulo en la ubicación o su creación") {
					itemLoc.Stock = int(datum.Quantity)
					s.storeStockRepo.Update(&stock)
					s.itemlocationrepo.Update(&itemLoc)
					//}
				}
			}
		}
	}
	// if err := s.storeRepo.DB.Exec("CALL ProcesarProductosAgrupados();").Error; err != nil {
	// 	log.Printf("Error ejecutando ProcesarProductosAgrupados: %v", err)
	// }
	return utils.ReturnUpdateErrorsExcel(stockErr), "ErrorsOnUpdate.xlsx", nil
}

type FatherData struct {
	FatherSku string
	FatherId  uint64
}

func returnFatherData(sku string) (FatherData, error) {
	var father FatherData
	var item2 models.Item
	var item *models.Item
	var err error
	itemRepo := itemsrepo.NewItemRepository(db.DB)
	parentRepo := itemsparentsrepo.NewItemParentRepository(db.DB)
	item2, err = itemRepo.FindByMainSku(sku)
	if err != nil {
		return father, err
	}
	if item2.ItemType == "son" {
		rel, err2 := parentRepo.FindByChild(item2.ID)
		if err2 != nil {
			return father, err2
		} else {
			item, _ = itemRepo.FindByID(rel.ParentItemID)
		}

	} else {
		item = &item2
	}
	father.FatherSku = item.MainSKU
	father.FatherId = item.ID
	return father, nil

}

// Obtiene los stock de un alamcén en base a su nombre
// A su vez el stock puede ser filtrado en base al parametro searchParam
func (s *StoreServiceImpl) GetStoreStock(storeName string, page int, pageSize int, searchParam string) []dtos.ItemStockInfo {
	store := s.storeRepo.FindByName(storeName)
	var results []dtos.ItemStockInfo
	var stock []models.StoreStock

	if searchParam == "" {
		stock, _ = s.storeStockRepo.FindByStore(store.ID, page, pageSize)
	} else {

		stock, _ = s.storeStockRepo.FindByStoreAndSearchParams(store.ID, searchParam, page, pageSize)

	}
	for _, itemInStock := range stock {
		childs, _ := s.itemsParentsRepo.FindByParent(itemInStock.Item.ID, 3, 0)
		results = append(results, dtos.ItemStockInfo{
			Itemname: itemInStock.Item.Name,
			Ean:      itemInStock.Item.EAN,
			SKU:      itemInStock.SKU_Parent,
			Amount:   itemInStock.Amount,
			Childs:   getChilds(childs),
		})
	}

	return results

}
