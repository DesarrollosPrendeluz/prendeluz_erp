package services

import (
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/fatherorderrepo"
	"prendeluz/erp/internal/repositories/itemlocationrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/stockdeficitrepo"
	"prendeluz/erp/internal/repositories/storestockrepo"
	"strconv"
)

type FatherOrderImpl struct {
	fatherorderrepo  fatherorderrepo.FatherOrderImpl
	orderrepo        orderrepo.OrderRepoImpl
	itemsRepo        itemsrepo.ItemRepoImpl
	orderitemrepo    orderitemrepo.OrderItemRepoImpl
	storestockrepo   storestockrepo.StoreStockRepoImpl
	itemlocationrepo itemlocationrepo.ItemLocationImpl
	stockdeficitrepo stockdeficitrepo.StockDeficitImpl
}

func NewFatherOrderService() *FatherOrderImpl {
	fatherorderrepo := *fatherorderrepo.NewFatherOrderRepository(db.DB)
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)
	orderitemrepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	orderrepo := *orderrepo.NewOrderRepository(db.DB)
	storestockrepo := *storestockrepo.NewStoreStockRepository(db.DB)
	itemlocationrepo := *itemlocationrepo.NewInItemLocationRepository(db.DB)
	stockdeficitrepo := *stockdeficitrepo.NewStockDeficitRepository(db.DB)
	return &FatherOrderImpl{
		fatherorderrepo:  fatherorderrepo,
		itemsRepo:        itemsRepo,
		orderitemrepo:    orderitemrepo,
		orderrepo:        orderrepo,
		storestockrepo:   storestockrepo,
		itemlocationrepo: itemlocationrepo,
		stockdeficitrepo: stockdeficitrepo}
}

func (s *FatherOrderImpl) FindLinesByFatherOrderCode(pageSize int, offset int, fatherOrderCode string, ean string, supplier_sku string, storeId int) (dtos.FatherOrderOrdersAndLines, int64, error) {
	var result dtos.FatherOrderOrdersAndLines
	var items []models.OrderItem
	var totalRecords int64
	var lines []dtos.LinesInfo
	var itemIds []uint64

	parentData, orderIds, _ := s.fatherorderrepo.FindParentAndOrders(fatherOrderCode)

	itemIds, _ = s.itemsRepo.FindByEanAndSupplierSku(ean, supplier_sku)

	items, totalRecords = s.orderitemrepo.FindByOrderAndItem(orderIds, storeId, itemIds, offset, pageSize)

	//procesado de datos de la query de lineas

	for _, item := range items {
		// Obtener el nombre del proveedor

		supplierName, supplierRef := returnSupplierData(item)
		locations := returnLocations(item)
		var fatherSku string

		if item.Item.ItemType == models.Father {
			fatherSku = item.Item.MainSKU
		} else {
			fatherSku = item.Item.FatherRel.Parent.MainSKU
		}

		// Crear la línea de información
		lineInfo := dtos.LinesInfo{
			LineID:          uint(item.ID),
			OrderCode:       item.OrderID,
			Name:            *item.Item.Name,
			Quantity:        int(item.Amount),
			RecivedQuantity: int(item.RecivedAmount),
			MainSku:         item.Item.MainSKU,
			Ean:             item.Item.EAN,
			FatherMainSku:   fatherSku,
			SupplierName:    supplierName,
			SupplierRef:     supplierRef,
			Location:        locations,
			AssignedUser: dtos.AssignedUserToOrderItem{
				AssignationId: item.AssignedRel.ID,
				UserId:        item.AssignedRel.UserID,
				UserName:      item.AssignedRel.UserRel.Name,
			},
		}

		// Añadir la línea al resultado
		lines = append(lines, lineInfo)
	}
	//Monatje de lineas
	result.FatherOrder = parentData
	result.Lines = lines

	return result, totalRecords, nil
}

func (s *FatherOrderImpl) CloseOrderByFather(fatherOrderId uint64) error {
	fatherData, _ := s.fatherorderrepo.FindByID(fatherOrderId)
	orderData, _ := s.orderrepo.FindByFatherId(fatherData.ID)
	for _, order := range orderData {
		linesData, _ := s.orderitemrepo.FindByOrder(order.ID)
		for _, line := range linesData {
			if line.RecivedAmount < line.Amount {
				var fatherSku string
				var location uint64
				diffAmount := line.Amount - line.RecivedAmount
				item, _ := s.itemsRepo.FindByIdWithFatherPreload(line.ItemID)

				if item.ItemType == "father" {
					fatherSku = item.MainSKU
				} else {
					fatherSku = item.FatherRel.Parent.MainSKU
				}

				switch line.StoreID {
				case 1:
					location = 1
				case 2:
					location = 86

				}

				line.RecivedAmount = line.Amount
				s.orderitemrepo.Update(&line)

				itemStock, _ := s.storestockrepo.FindByItemAndStore(fatherSku, strconv.FormatInt(line.StoreID, 10))
				itemStock.Amount = itemStock.Amount + diffAmount
				s.storestockrepo.Update(&itemStock)

				itemStockLocation, _ := s.itemlocationrepo.FindByItemsAndLocation(fatherSku, location)
				itemStockLocation.Stock = itemStockLocation.Stock + int(diffAmount)
				s.itemlocationrepo.Update(&itemStockLocation)

				stockDef, _ := s.stockdeficitrepo.GetByFatherAndStore(fatherSku, line.StoreID)

				stockDef.Amount = stockDef.Amount - diffAmount
				stockDef.PendingAmount = stockDef.PendingAmount - diffAmount
				if stockDef.Amount < 0 {
					stockDef.Amount = 0
				}
				if stockDef.PendingAmount < 0 {
					stockDef.Amount = 0
				}
				s.stockdeficitrepo.Update(&stockDef)

			}
		}

		order.OrderStatusID = 3
		s.orderrepo.Update(&order)

	}

	fatherData.OrderStatusID = 3
	s.fatherorderrepo.Update(fatherData)

	//s.stockdeficitrepo.CallStockDefProc()
	//s.stockdeficitrepo.CallPendingStockProc()

	return nil

}

func returnLocations(item models.OrderItem) []string {
	var locs []string
	var locations *[]models.ItemLocation

	// Determinar qué lista de ubicaciones usar
	if item.Item.ItemType != models.Father && item.Item.FatherRel != nil && item.Item.FatherRel.Parent != nil {
		locations = item.Item.FatherRel.Parent.ItemLocations
	} else {
		locations = item.Item.ItemLocations
	}

	// Recorrer y agregar ubicaciones si existen
	if locations != nil && len(*locations) > 0 {
		for _, location := range *locations {
			locs = append(locs, location.StoreLocations.Name)
		}
	} else {
		locs = append(locs, "")
	}

	return locs
}

func returnSupplierData(item models.OrderItem) (string, string) {
	supplierName, supplierRef := "", ""

	var supplierItems *[]models.SupplierItem

	if item.Item.ItemType != models.Father {
		if item.Item.FatherRel != nil && item.Item.FatherRel.Parent != nil {
			supplierItems = item.Item.FatherRel.Parent.SupplierItems
		}
	} else {
		supplierItems = item.Item.SupplierItems
	}

	if supplierItems != nil && len(*supplierItems) > 0 {
		firstSupplierItem := (*supplierItems)[0]
		if firstSupplierItem.Supplier != nil {
			supplierName = firstSupplierItem.Supplier.Name
			supplierRef = firstSupplierItem.SupplierSku
		}
	}

	return supplierName, supplierRef
}
