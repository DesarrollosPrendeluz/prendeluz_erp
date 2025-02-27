package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/asinrepo"
	"prendeluz/erp/internal/repositories/boxrepo"
	"prendeluz/erp/internal/repositories/fatherorderrepo"
	"prendeluz/erp/internal/repositories/itemlocationrepo"
	"prendeluz/erp/internal/repositories/itemsparentsrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderlineboxrepo"
	"prendeluz/erp/internal/repositories/orderlinelocationviewrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/outorderrelationrepo"
	"prendeluz/erp/internal/repositories/palletrepo"
	"prendeluz/erp/internal/repositories/stockdeficitrepo"
	"prendeluz/erp/internal/repositories/storestockrepo"
	"prendeluz/erp/internal/repositories/supplieritemrepo"
	"prendeluz/erp/internal/repositories/supplierorderrepo"
	"prendeluz/erp/internal/repositories/suppliersoldorderrelationrepo"
	stockservices "prendeluz/erp/internal/services/stock_deficit"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelExportation struct {
	OC_code string
	Asin    string
	Total   int
	Per_box float64
	Rest    string
	Pallet  string
	Box     string
}

type FatherOrderImpl struct {
	orderlinelocationviewrepo     orderlinelocationviewrepo.OrderLineLocationViewImpl
	supplierorderrepo             supplierorderrepo.SupplierOrderImpl
	fatherorderrepo               fatherorderrepo.FatherOrderImpl
	orderrepo                     orderrepo.OrderRepoImpl
	itemsRepo                     itemsrepo.ItemRepoImpl
	orderitemrepo                 orderitemrepo.OrderItemRepoImpl
	storestockrepo                storestockrepo.StoreStockRepoImpl
	itemlocationrepo              itemlocationrepo.ItemLocationImpl
	stockdeficitrepo              stockdeficitrepo.StockDeficitImpl
	asinrepo                      asinrepo.AsinRepoImpl
	itemsparentsrepo              itemsparentsrepo.ItemsParentsRepoImpl
	palletsrepo                   palletrepo.PalletImpl
	boxesrepo                     boxrepo.BoxImpl
	orderlineboxrepo              orderlineboxrepo.OrderLineBoxImpl
	suppliersoldorderrelationrepo suppliersoldorderrelationrepo.SupplierSoldOrderRelationImpl
}

func NewFatherOrderService() *FatherOrderImpl {
	fatherorderrepo := *fatherorderrepo.NewFatherOrderRepository(db.DB)
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)
	orderitemrepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	orderrepo := *orderrepo.NewOrderRepository(db.DB)
	storestockrepo := *storestockrepo.NewStoreStockRepository(db.DB)
	itemlocationrepo := *itemlocationrepo.NewInItemLocationRepository(db.DB)
	stockdeficitrepo := *stockdeficitrepo.NewStockDeficitRepository(db.DB)
	asinrepo := *asinrepo.NewAsinRepository(db.DB)
	supplierorderrepo := *supplierorderrepo.NewSupplierOrderRepository(db.DB)
	itemsparentsrepo := *itemsparentsrepo.NewItemParentRepository(db.DB)
	palletsrepo := *palletrepo.NewPalletRepository(db.DB)
	boxesrepo := *boxrepo.NewBoxRepository(db.DB)
	orderlineboxrepo := *orderlineboxrepo.NewOrderLineBoxRepository(db.DB)
	orderlinelocationviewrepo := *orderlinelocationviewrepo.NewOrderLineLocationViewRepository(db.DB)
	suppliersoldorderrelationrepo := *suppliersoldorderrelationrepo.NewSupplierSoldOrderRelationRepository(db.DB)

	return &FatherOrderImpl{
		fatherorderrepo:               fatherorderrepo,
		itemsRepo:                     itemsRepo,
		orderitemrepo:                 orderitemrepo,
		orderrepo:                     orderrepo,
		storestockrepo:                storestockrepo,
		itemlocationrepo:              itemlocationrepo,
		stockdeficitrepo:              stockdeficitrepo,
		asinrepo:                      asinrepo,
		supplierorderrepo:             supplierorderrepo,
		orderlinelocationviewrepo:     orderlinelocationviewrepo,
		itemsparentsrepo:              itemsparentsrepo,
		palletsrepo:                   palletsrepo,
		boxesrepo:                     boxesrepo,
		orderlineboxrepo:              orderlineboxrepo,
		suppliersoldorderrelationrepo: suppliersoldorderrelationrepo,
	}

}

func (s *FatherOrderImpl) FindLinesByFatherOrderCode(pageSize int, offset int, fatherOrderCode string, ean string, supplier_sku string, storeId int, searchByEan string, searchByLoc string, locFilter string) (dtos.FatherOrderOrdersAndLines, int64, error) {
	var result dtos.FatherOrderOrdersAndLines
	var items []models.OrderItem
	var totalRecords int64
	var lines []dtos.LinesInfo
	var itemIds []uint64
	calcPage := offset * pageSize

	parentData, orderIds, _ := s.fatherorderrepo.FindParentAndOrders(fatherOrderCode)
	if locFilter != "" {
		list, _ := s.orderlinelocationviewrepo.FindLineArrayByFatherAndLocation(parentData.ID, storeId, locFilter)
		items, totalRecords = s.orderitemrepo.FindByLineID(list, calcPage, pageSize)

	} else if searchByEan != "" && searchByLoc != "" && ean == "" {
		list, order, _ := s.orderlinelocationviewrepo.FindByFatherAndStoreWithOrder(parentData.ID, storeId, searchByLoc, searchByEan)
		items, totalRecords = s.orderitemrepo.FindByLineIDWithOrder(list, order, calcPage, pageSize)

	} else {
		itemIds, _ = s.itemsRepo.FindByEanAndSupplierSku(ean, supplier_sku)
		items, totalRecords = s.orderitemrepo.FindByOrderAndItem(orderIds, storeId, itemIds, calcPage, pageSize)
	}

	//procesado de datos de la query de lineas

	for _, item := range items {
		// Obtener el nombre del proveedor

		supplierName, supplierRef := returnSupplierData(item, parentData.GenericSupplier)
		locations := returnLocations(item, uint64(storeId))
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
			Box:             item.Box,
			Pallet:          item.Pallet,
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
func (s *FatherOrderImpl) ClosePickingOrders(fatherOrderId uint64) error {
	orders, error := s.orderrepo.FindByFatherId(fatherOrderId)
	if error == nil {
		for _, order := range orders {
			orderLines, errl := s.orderitemrepo.FindByOrder(order.ID)
			if errl == nil {
				for _, orderLine := range orderLines {
					if orderLine.StoreID == 1 {
						item, _ := s.itemsRepo.FindByID(orderLine.ItemID)
						if item.ItemType == "son" {
							fatherRel, _ := s.itemsparentsrepo.FindByChild(item.ID)
							item, _ = s.itemsRepo.FindByID(fatherRel.ParentItemID)
						}
						locations, _ := s.itemlocationrepo.FindByItemsAndStore(item.MainSKU, 1, -1, -1)
						stockToRest := orderLine.Amount
						for _, location := range locations {
							if stockToRest > 0 && location.Stock > 0 {
								if stockToRest <= int64(location.Stock) {
									location.Stock = (location.Stock - int(stockToRest))
									stockToRest = 0

								} else {
									location.Stock = 0
									stockToRest = (stockToRest - int64(location.Stock))

								}
								s.itemlocationrepo.Update(&location)
							}

						}

						stock, _ := s.storestockrepo.FindByItemAndStore(item.MainSKU, "1")
						stock.Amount = (stock.Amount - orderLine.Amount)
						s.storestockrepo.Update(&stock)

						orderLine.RecivedAmount = orderLine.Amount
						s.orderitemrepo.Update(&orderLine)

					}

				}
			}

		}

	}

	return nil

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

func (s *FatherOrderImpl) OpenOrderByFather(fatherOrderId uint64) error {
	fatherData, _ := s.fatherorderrepo.FindByID(fatherOrderId)
	orderData, _ := s.orderrepo.FindByFatherId(fatherData.ID)
	for _, order := range orderData {
		linesData, _ := s.orderitemrepo.FindByOrder(order.ID)
		for _, line := range linesData {
			if line.RecivedAmount == line.Amount {
				var fatherSku string
				var location uint64
				diffAmount := line.Amount
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

				line.RecivedAmount = 0
				s.orderitemrepo.Update(&line)
				/*Quitamos el stock del total*/
				itemStock, _ := s.storestockrepo.FindByItemAndStore(fatherSku, strconv.FormatInt(line.StoreID, 10))
				itemStock.Amount = itemStock.Amount - diffAmount
				s.storestockrepo.Update(&itemStock)
				/*Quitamos el stock de la ubicación por defecto*/
				itemStockLocation, _ := s.itemlocationrepo.FindByItemsAndLocation(fatherSku, location)
				itemStockLocation.Stock = itemStockLocation.Stock - int(diffAmount)
				s.itemlocationrepo.Update(&itemStockLocation)
				/*Actualizamos el stock deficit*/
				stockDef, _ := s.stockdeficitrepo.GetByFatherAndStore(fatherSku, line.StoreID)

				stockDef.Amount = stockDef.Amount + diffAmount
				stockDef.PendingAmount = stockDef.PendingAmount + diffAmount
				if stockDef.Amount < 0 {
					stockDef.Amount = 0
				}
				if stockDef.PendingAmount < 0 {
					stockDef.Amount = 0
				}
				s.stockdeficitrepo.Update(&stockDef)

			}
		}

		order.OrderStatusID = 1
		s.orderrepo.Update(&order)

	}

	fatherData.OrderStatusID = 1
	s.fatherorderrepo.Update(fatherData)

	//s.stockdeficitrepo.CallStockDefProc()
	//s.stockdeficitrepo.CallPendingStockProc()

	return nil

}

func returnLocations(item models.OrderItem, store_id uint64) []string {
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
		sort.Slice(*locations, func(i, j int) bool {
			return (*locations)[i].Stock > (*locations)[j].Stock
		})

		for _, location := range *locations {
			if location.StoreLocations.StoreID == store_id {
				locs = append(locs, location.StoreLocations.Code)
			}

		}
	} else {
		locs = append(locs, "")
	}

	return locs
}

func returnSupplierData(item models.OrderItem, supplier *models.SupplierOrder) (string, string) {
	supplierName, supplierRef := "", ""

	var supplierItems *[]models.SupplierItem
	if supplier != nil && supplier.SupplierID != 0 {
		var id uint64
		if item.Item.ItemType != models.Father {
			id = item.Item.FatherRel.Parent.ID
		} else {
			id = item.Item.ID
		}
		item, _ := supplieritemrepo.NewSupplierItemRepository(db.DB).FindBySupplierIdAndItemId(id, supplier.SupplierID)
		return supplier.Supplier.Name, item.SupplierSku

	} else {
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

	}

	return supplierName, supplierRef
}
func (s *FatherOrderImpl) DownloadExcelAmazon(fatherID uint64) string {
	father, _ := s.fatherorderrepo.FindByID(fatherID)
	_, childOrders, _ := s.fatherorderrepo.FindParentAndOrders(father.Code)
	var results []ExcelExportation
	for _, orderId := range childOrders {
		order, _ := s.orderrepo.FindByID(orderId)

		lines, _ := s.orderitemrepo.FindByOrder(order.ID)
		for _, itemline := range lines {
			asin, _ := s.asinrepo.FindByItemId(uint64(itemline.ItemID))

			boxlines, _ := s.orderlineboxrepo.GetByLineId(int(itemline.ID)) //POr linea es en realidad
			if len(boxlines) != 0 {
				fmt.Println(boxlines)
			}
			for _, boxline := range boxlines {
				boxNumber, _ := s.boxesrepo.FindByID(uint64(boxline.BoxID))

				if len(boxlines) != 0 {
					fmt.Println(boxNumber)

				}
				palletNUmber, _ := s.palletsrepo.FindByID(boxNumber.PalletID)
				tmp := ExcelExportation{
					OC_code: order.Code,
					Asin:    asin.Code,
					Box:     strconv.Itoa(boxNumber.Number),
					Pallet:  strconv.Itoa(palletNUmber.Number),
					Per_box: float64(boxline.Quantity),
					Total:   int(itemline.Amount),
				}
				results = append(results, tmp)
			}
		}
	}
	return generateExcelBase64(results)

}

// DEPRECATED
func (s *FatherOrderImpl) DownloadOrdersExcelToAmazon(fatherID uint64) string {
	var exportData []ExcelExportation
	var boxSubStrings []string
	var palletSubStrings []string

	subStringsDivider := func(data *string) []string {
		if data != nil && *data != "" {
			return strings.Split(*data, ",")
		}
		return []string{"-"}
	}

	fatherData, fatherError := s.orderrepo.FindByFatherId(fatherID)
	if fatherError != nil {
		fmt.Println(fatherError.Error())
	}
	for _, father := range fatherData {
		orderItems, orderItemError := s.orderitemrepo.FindByOrderAndStore(father.ID, 2)
		if orderItemError != nil {
			fmt.Println(orderItemError.Error())
		}
		for _, orderItem := range orderItems {
			asin, asinError := s.asinrepo.FindByItemId(orderItem.ItemID)
			// var boxSubStrings []string
			if asinError != nil {
				fmt.Println(fatherError.Error())
			}
			boxSubStrings = subStringsDivider(orderItem.Box)
			palletSubStrings = subStringsDivider(orderItem.Pallet)
			numberOfBoxes := len(boxSubStrings) * len(palletSubStrings)
			for _, pallet := range palletSubStrings {
				for _, boxIndv := range boxSubStrings {
					partials := float64(orderItem.RecivedAmount) / float64(numberOfBoxes)
					data := ExcelExportation{
						OC_code: father.Code,
						Asin:    asin.Code,
						Total:   int(orderItem.RecivedAmount),
						Per_box: partials,
						Box:     boxIndv,
						Pallet:  pallet,
					}
					exportData = append(exportData, data)

				}

			}

		}

	}

	return generateExcelBase64(exportData)

}

func (s *FatherOrderImpl) CreateOrder(requestBody dtos.OrderWithLinesRequest) bool {

	var code string

	fechaActual := time.Now().Format("2006-01-02 15:04:05")
	code = "request.generated." + fechaActual

	// Acceder a los valores del cuerpo
	for _, dataItem := range requestBody.Data {
		order := dataItem.Order
		lines := dataItem.Lines
		fatherRepo := s.fatherorderrepo
		repo := s.orderrepo
		if order.Name != nil {
			code = *order.Name
		}
		fatherObject := models.FatherOrder{
			OrderStatusID: order.Status,
			OrderTypeID:   order.Type,
			Code:          code,
			Filename:      "request",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if fatherRepo.Create(&fatherObject) == nil {
			orderObject := models.Order{
				OrderStatusID: order.Status,
				FatherOrderID: fatherObject.ID,
				Code:          "request.generated." + fechaActual,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			if repo.Create(&orderObject) == nil {
				createOrderLines(fatherObject, orderObject, lines)

			}

			if dataItem.Order.Supplier != nil {
				supplierOrderObject := models.SupplierOrder{
					SupplierID:    *dataItem.Order.Supplier,
					FatherOrderID: fatherObject.ID,
				}
				s.supplierorderrepo.Create(&supplierOrderObject)
			}
			if fatherObject.OrderTypeID == uint64(1) {
				fatherRelItem, _ := s.fatherorderrepo.FindLatestByType(2)
				s.suppliersoldorderrelationrepo.Create(&models.SupplierSoldOrderRelation{
					SupplierID:  uint(fatherObject.ID),
					SoldOrderID: uint(fatherRelItem.ID),
				})

			}

		}
		//README: Por el funcionamiento de la aplicación se ha decidido no ejecutar los procedimientos almacenados
		// if err := db.DB.Exec("CALL UpdateStockDeficitByStore();").Error; err != nil {
		// 	log.Printf("Error ejecutando UpdateStockDeficitByStore: %v", err)
		// } else {
		// 	fmt.Println("en teoría se ha ejecutado: CALL UpdateStockDeficitByStore();")

		// }

		// // Llamada al segundo procedimiento almacenado
		// if err := db.DB.Exec("CALL UpdatePendingStocks();").Error; err != nil {
		// 	log.Printf("Error ejecutando UpdatePendingStocks: %v", err)
		// } else {
		// 	fmt.Println("en teoría se ha ejecutado: CALL UpdatePendingStocks()")

		// }

	}
	return true

}

func (s *FatherOrderImpl) FindAllWithAssocData(fatherOrderCode string, typeId int, statusId int, pageSize int, offset int) ([]dtos.FatherOrderWithRecount, int64, error) {
	calcPage := pageSize * offset
	return s.fatherorderrepo.FindAllWithAssocData(pageSize, calcPage, fatherOrderCode, typeId, statusId)

}

func createOrderLines(fatherOrder models.FatherOrder, order models.Order, lines []dtos.Line) error {
	repo := orderitemrepo.NewOrderItemRepository(db.DB) // Asumiendo que tienes un repositorio para las líneas
	//itemRepo := itemsrepo.NewItemRepository(db.DB)

	for _, line := range lines {
		//son, _ := itemRepo.FindSonId(line.ItemID)

		orderLine := models.OrderItem{
			OrderID:       order.ID,
			ItemID:        line.ItemID,
			Amount:        line.Quantity,
			RecivedAmount: line.RecivedQuantity,
			StoreID:       line.StoreID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Guardar cada línea en la base de datos
		if err := repo.Create(&orderLine); err != nil {

			return err
		}
		if fatherOrder.OrderTypeID == uint64(2) && line.ClientID != nil {
			outRelRepo := outorderrelationrepo.NewOutOrderRelationRepository(db.DB)
			outRel := models.OutOrderRelation{
				ClientID:    *line.ClientID,
				OrderLineID: orderLine.ID,
			}
			outRelRepo.Create(&outRel)

		}
		if fatherOrder.OrderTypeID == uint64(1) {
			stockservices.NewStockDeficitService().AddPendingStockByItem(line.ItemID, line.StoreID, int(line.Quantity))
		}
	}
	return nil
}

func generateExcelBase64(exportData []ExcelExportation) string {
	f := excelize.NewFile()
	// Crear encabezados en la primera fila
	sheetName := "Amazon_Data"

	f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A1", "PO")
	f.SetCellValue(sheetName, "B1", "Asin")
	f.SetCellValue(sheetName, "C1", "Total")
	f.SetCellValue(sheetName, "D1", "Per box")
	f.SetCellValue(sheetName, "E1", "Pallet")
	f.SetCellValue(sheetName, "F1", "Box")
	f.SetCellValue(sheetName, "G1", "Pallet Label")
	f.SetCellValue(sheetName, "H1", "Box Label")

	// Escribir los datos en las filas siguientes
	for i, data := range exportData {
		row := i + 2 // La primera fila es para los encabezados

		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), data.OC_code)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), data.Asin)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), data.Total)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), data.Per_box)
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), data.Pallet)
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), data.Box)

	}
	f.DeleteSheet("Sheet1")
	// Escribir el archivo Excel en un buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return ""
	}

	// Codificar el contenido del buffer en Base64
	return base64.StdEncoding.EncodeToString(buf.Bytes())

}
