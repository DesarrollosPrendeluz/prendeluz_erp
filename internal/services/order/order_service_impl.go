package services

import (
	"fmt"
	"io"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"prendeluz/erp/internal/repositories/erpupdateorderlinehistoryrepo"
	"prendeluz/erp/internal/repositories/fatherorderrepo"
	"prendeluz/erp/internal/repositories/itemsparentsrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/ordererrorrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/outorderrelationrepo"
	"prendeluz/erp/internal/repositories/stockdeficitrepo"
	stockrepo "prendeluz/erp/internal/repositories/storestockrepo"
	"prendeluz/erp/internal/repositories/suppliersoldorderrelationrepo"
	"prendeluz/erp/internal/repositories/tokenrepo"
	stockDeficit "prendeluz/erp/internal/services/stock_deficit"
	"prendeluz/erp/internal/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ParentItemResult struct {
	ParentItemID int    `gorm:"column:parent_item_id"`
	MainSKU      string `gorm:"column:main_sku"`
}

type OrderServiceImpl struct {
	orderRepo                     orderrepo.OrderRepoImpl
	orderItemsRepo                orderitemrepo.OrderItemRepoImpl
	fatherOrderRepo               fatherorderrepo.FatherOrderImpl
	orderErrorRepo                repositories.GORMRepository[models.ErrorOrder]
	itemsRepo                     itemsrepo.ItemRepoImpl
	stockdeficitrepo              stockdeficitrepo.StockDeficitImpl
	erpupdateorderlinehistoryrepo erpupdateorderlinehistoryrepo.ErpUpdateOrderLineHistoryImpl
}

func NewOrderService() *OrderServiceImpl {
	orderRepo := *orderrepo.NewOrderRepository(db.DB)
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})
	orderItemRepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)
	fatherOrderRepo := *fatherorderrepo.NewFatherOrderRepository(db.DB)
	stockdeficitrepo := *stockdeficitrepo.NewStockDeficitRepository(db.DB)
	erpupdateorderlinehistoryrepo := *erpupdateorderlinehistoryrepo.NewErpUpdateOrderLineHistoryRepository(db.DB)

	return &OrderServiceImpl{
		orderRepo:                     orderRepo,
		orderItemsRepo:                orderItemRepo,
		orderErrorRepo:                errorOrderRepo,
		itemsRepo:                     itemsRepo,
		fatherOrderRepo:               fatherOrderRepo,
		stockdeficitrepo:              stockdeficitrepo,
		erpupdateorderlinehistoryrepo: erpupdateorderlinehistoryrepo}
}

// Carga el excel y crea las nuevas ordenes en este caso solo de ventas por el momento
func (s *OrderServiceImpl) UploadOrderExcel(file io.Reader, filename string) error {

	fatherRepo := fatherorderrepo.NewFatherOrderRepository(db.DB)

	excelOrderList, err := utils.ExceltoJSON(file)

	if err != nil {
		return err
	}
	// for _, order := range excelOrderList {
	// 	fmt.Printf("Contenido del orden: %+v\n", order)
	// }
	//pendiente de crear la father order
	fatherObject := models.FatherOrder{
		OrderStatusID: uint64(orderrepo.Order_Status["pediente"]),
		OrderTypeID:   uint64(orderrepo.Order_Types["venta"]),
		Code:          quitarExtension(filename),
		Filename:      filename,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if fatherRepo.Create(&fatherObject) == nil {

		succesOrders, orderItems, errorOrders := generateOrdersAndOrderLines(excelOrderList, fatherObject.ID)
		orderRepo := orderrepo.NewOrderRepository(db.DB)
		orderItem := orderitemrepo.NewOrderItemRepository(db.DB)
		out := outorderrelationrepo.NewOutOrderRelationRepository(db.DB)
		orderErr := ordererrorrepo.NewOrderErrRepository(db.DB)

		_, errs := orderRepo.CreateAll(&succesOrders)
		if errs != nil {
			return err
		}

		if len(errorOrders) > 0 {
			_, err := orderErr.CreateAll(&errorOrders)
			if err != nil {
				return err
			}
		}
		var orderItemToInsert []models.OrderItem
		for _, order := range succesOrders {
			for _, tmp := range orderItems[order.Code] {
				tmp.OrderID = order.ID
				orderItemToInsert = append(orderItemToInsert, tmp)
			}
		}
		for _, orderLine := range orderItemToInsert {
			err = orderItem.Create(&orderLine)
			out.Create(&models.OutOrderRelation{
				ClientID:    orderLine.ClientID,
				OrderLineID: orderLine.ID,
			})

		}
		stockDeficit.NewStockDeficitService().CalcStockDeficitByFatherOrder(fatherObject.ID)

		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}

}

// Obtiene las ordenes paginadas en base a los parámetros page y pagesize
// A su vez si recibe los parámetros start dat y end date se filtran dichas ordenes por fecha de creación
func (s *OrderServiceImpl) GetOrders(page int, pageSize int, startDate string, endDate string, statusId int, orderCode string) ([]dtos.ItemsPerOrder, int64, error) {

	var results []dtos.ItemsPerOrder
	var orders []models.Order
	var recount int64
	var err error
	offset := (page - 1) * pageSize

	if startDate != "" || endDate != "" || orderCode != "" || statusId != 0 {
		orders, recount, err = s.orderRepo.FindOrderFiltered(pageSize, offset, startDate, endDate, statusId, orderCode)

	} else {
		orders, recount, err = s.orderRepo.FindAll(pageSize, offset)
	}

	if err != nil {
		return results, recount, err
	}

	for _, order := range orders {
		var itemOrder dtos.ItemsPerOrder
		fmt.Println(order.ID)
		itemOrder.Id = order.ID
		itemOrder.OrderCode = order.Code
		itemOrder.TypeID = int64(order.FatherOrder.OrderTypeID)
		itemOrder.StatusID = int64(order.OrderStatusID)
		itemOrder.Type = order.FatherOrder.OrderType.Name
		itemOrder.Status = order.OrderStatus.Name
		orderItemList, _ := s.orderItemsRepo.FindByOrder(order.ID)

		for _, orderItem := range orderItemList {
			var itemInfo dtos.ItemInfo
			item, _ := s.itemsRepo.FindByIdExtraData(orderItem.ItemID)
			itemInfo.Id = orderItem.ID
			itemInfo.Sku = item.MainSKU
			itemInfo.Name = *item.Name
			itemInfo.Ean = item.EAN
			itemInfo.Amount = orderItem.Amount
			itemInfo.RecivedAmount = orderItem.RecivedAmount
			//FIXME: quitar una vez se haya externalizado
			itemInfo.Box = *orderItem.Box
			itemInfo.Pallet = *orderItem.Pallet
			if item.SupplierItems != nil && len(*item.SupplierItems) > 0 && (*item.SupplierItems)[0].Supplier != nil {
				itemInfo.Supplier = (*item.SupplierItems)[0].Supplier.Name
			} else {
				itemInfo.Supplier = "No asignado"
			}
			if item.ItemLocations != nil && len(*item.ItemLocations) > 0 {
				var temLocation []string
				for _, location := range *item.ItemLocations {
					if location.StoreLocations != nil {
						temLocation = append(temLocation, location.StoreLocations.Code)
					}
				}
				itemInfo.Locations = temLocation

			}

			if orderItem.AssignedRel.ID != 0 {
				itemInfo.AssignedUser.AssignationId = orderItem.AssignedRel.ID
				itemInfo.AssignedUser.UserId = uint64(orderItem.AssignedRel.UserRel.ID)
				itemInfo.AssignedUser.UserName = orderItem.AssignedRel.UserRel.Name
			}
			itemOrder.ItemsOrdered = append(itemOrder.ItemsOrdered, itemInfo)

		}
		results = append(results, itemOrder)

	}

	return results, recount, nil

}

// Actualiza el estado de una orden a completada
func (s *OrderServiceImpl) OrderComplete(orderCode string) error {
	order, err := s.orderRepo.FindByOrderCode(orderCode)
	if err != nil {
		return err
	}

	s.orderRepo.UpdateStatus(orderrepo.Order_Status["finalizada"], order.ID)

	return nil
}

// Carga el excel y crea las nuevas ordenes en este caso solo de ventas por el momento
func (s *OrderServiceImpl) UploadOrdersByExcel(file io.Reader, requestFatherOrderCode string, token string) (string, string) {
	var orderIdArr []uint64
	var addErrData []utils.UpdateOrderError
	addError := func(errorData error, errArr *[]utils.UpdateOrderError, sku string, err string) bool {
		if errorData != nil {
			errReturn := utils.UpdateOrderError{
				FatherSku: sku,
				Error:     err,
			}
			*errArr = append(*errArr, errReturn)
			return false

		}
		return true

	}
	//asiganmos el padre de la orden si lo hay y el order id
	if requestFatherOrderCode != "" {
		fatherOrder, fatherError := s.fatherOrderRepo.FindByCode(requestFatherOrderCode)
		if addError(fatherError, &addErrData, "", "No se ha encontrado la order padre") {
			orders, orderError := s.orderRepo.FindByFatherId(fatherOrder.ID)
			excelOrderList, _ := utils.ExcelToJSONOrder(file)
			if addError(orderError, &addErrData, "", "No se han  encontrado las ordenes hijas") {

				for _, order := range orders {
					orderIdArr = append(orderIdArr, order.ID)
				}
				currentDate := time.Now().Format("20060102")
				code := utils.GenerateRandomString(10) + "-" + currentDate

				for _, line := range excelOrderList {
					sku := line.Sku
					sku = strings.ReplaceAll(strings.ReplaceAll(sku, " ", ""), "\n", "") // Quitar espacios y saltos de línea
					items, _ := s.itemsRepo.FindByMainSku(sku)
					orderLine, orderLineError := s.orderItemsRepo.FindByItemAndOrders(orderIdArr, items.ID, 2)
					if addError(orderLineError, &addErrData, sku, "No se ha encontrado la linea del articulo") {
						updatesDeficitsByLine(items.ID, fatherOrder.OrderTypeID, orderLine.OrderID, line.Quantity, orderLine.Amount)
						ol := orderLine
						orderLine.Amount = line.Quantity
						repo := tokenrepo.NewTokenRepository(db.DB)
						user, _ := repo.ReturnDataByToken(token)
						s.erpupdateorderlinehistoryrepo.GenerateOrderLineHistory(ol, orderLine, user.UserId, line.Type, code)
						if fatherOrder.OrderTypeID == 1 {
							updateRelatedOrderProcess(fatherOrder.ID, items.EAN, orderLine.Amount, ol.Amount, user.UserId, line.Type, code)
						}
						s.orderItemsRepo.Update(&orderLine)
					}
					//proveedor

				}
				//update order status
				for _, order := range orders {
					order.OrderStatusID = uint64(orderrepo.Order_Status["en_espera"])
					s.orderRepo.Update(&order)
				}

				//update father order status
				fatherOrder.OrderStatusID = uint64(orderrepo.Order_Status["en_espera"])
				s.fatherOrderRepo.Update(&fatherOrder)
			}
		}

	} else {
		addError(fmt.Errorf("no se ha enviado orden padre"), &addErrData, "", "No se ha pasado la orden padre por parámetro")

	}
	return utils.ReturnUpdateOrdersErrorsExcel(addErrData), "update_errors.xlsx"

}

func updatesDeficitsByLine(itemId uint64, fatherOrderType uint64, orderId uint64, updateLineQuantity int64, orderLineQuantity int64) {
	parent := returnParentItemById(itemId)
	deficitRepo := stockdeficitrepo.NewStockDeficitRepository(db.DB)
	deficit, _ := deficitRepo.FindOrCreateByFatherAndStore(parent.MainSKU, 2)
	//proveedor
	if fatherOrderType == 1 {
		deficit.PendingAmount = (deficit.PendingAmount - orderLineQuantity) + updateLineQuantity

	} else if fatherOrderType == 2 {
		out := 0
		in := 0
		orderLines, _ := orderitemrepo.NewOrderItemRepository(db.DB).FindByItemOrderStore(itemId, orderId, 1)
		in = int(orderLines.Amount)
		orderLines, _ = orderitemrepo.NewOrderItemRepository(db.DB).FindByItemOrderStore(itemId, orderId, 2)
		out = int(orderLines.Amount)

		actualDiff := out - in
		futureDiff := int(updateLineQuantity) - in

		deficit.Amount = (deficit.Amount - int64(actualDiff)) + int64(futureDiff)

	}
	deficitRepo.Update(&deficit)

}

// retorna datos para crear ordenes las línea de las ordenes y los errores correspondientes
func generateOrdersAndOrderLines(rawOrders []utils.ExcelOrder, fatherOrderId uint64) ([]models.Order, map[string][]models.OrderItem, []models.ErrorOrder) {
	itemsRepo := itemsrepo.NewItemRepository(db.DB)

	var errorOrdersList []models.ErrorOrder
	var ordersList []models.Order
	orderItemsOk := make(map[string][]models.OrderItem)

	for _, orderCode := range rawOrders {
		var mainSkus []string
		for _, orderInfo := range orderCode.Info {
			mainSkus = append(mainSkus, orderInfo.MainSku)
		}
		itemsMap, _ := itemsRepo.FindByMainSkus(mainSkus)

		for _, orderInfo := range orderCode.Info {
			processOrderItem(orderInfo, itemsMap, orderCode.OrderCode, orderItemsOk, &errorOrdersList)
		}

		order := models.Order{
			OrderStatusID: uint64(orderrepo.Order_Status["iniciada"]),
			Code:          orderCode.OrderCode,
			FatherOrderID: fatherOrderId,
		}
		ordersList = append(ordersList, order)

	}
	return ordersList, orderItemsOk, errorOrdersList
}

func processOrderItem(orderInfo utils.OrderInfo, itemsMap map[string]models.Item, orderCode string, orderItemsOk map[string][]models.OrderItem, errorOrdersList *[]models.ErrorOrder) {
	item, exists := itemsMap[orderInfo.MainSku]

	if !exists {
		fmt.Println("error en el sku " + orderInfo.MainSku)
		errorOrder := models.ErrorOrder{
			Main_Sku: orderInfo.MainSku,
			Error:    "Item with sku " + orderInfo.MainSku + " not found",
			Order:    orderCode,
		}

		*errorOrdersList = append(*errorOrdersList, errorOrder)
	} else {
		orderItem := models.OrderItem{
			ItemID:        item.ID,
			Amount:        orderInfo.Amount,
			RecivedAmount: 0,
			StoreID:       int64(orderInfo.Store),
			ClientID:      orderInfo.Client,
		}
		orderItemsOk[orderCode] = append(orderItemsOk[orderCode], orderItem)
		addPickingLines(orderInfo, orderCode, item.ID, orderItemsOk, errorOrdersList)
	}
}

func addPickingLines(orderItemInfo utils.OrderInfo, orderCode string, itemId uint64, orderItemsOk map[string][]models.OrderItem, errorOrdersList *[]models.ErrorOrder) {
	stockrepo := stockrepo.NewStoreStockRepository(db.DB)
	parentStock, errorInParentStock := stockrepo.FindByItemAndStore(orderItemInfo.ParentSku, "1")

	if errorInParentStock != nil {
		fmt.Println("error en el sku " + orderItemInfo.ParentSku)
		errorOrder := models.ErrorOrder{
			Main_Sku: orderItemInfo.ParentSku,
			Error:    "Item with sku " + orderItemInfo.ParentSku + " not found" + errorInParentStock.Error(),
			Order:    orderCode,
		}

		*errorOrdersList = append(*errorOrdersList, errorOrder)
	} else {
		actualStock := parentStock.Amount - parentStock.ReservedAmount

		if parentStock.Amount > 0 && actualStock > 0 {

			orderItem := models.OrderItem{
				ItemID:        itemId,
				Amount:        parentStock.ReservedAmount,
				RecivedAmount: 0,
				StoreID:       1,
				ClientID:      orderItemInfo.Client,
			}
			if actualStock > orderItemInfo.Amount {
				parentStock.ReservedAmount += orderItemInfo.Amount
				stockrepo.Update(&parentStock)
				orderItem.Amount = orderItemInfo.Amount
			} else {
				parentStock.ReservedAmount += actualStock
				stockrepo.Update(&parentStock)
				orderItem.Amount = actualStock
			}

			if orderItem.Amount > 0 {
				orderItemsOk[orderCode] = append(orderItemsOk[orderCode], orderItem)
			}

		}

	}

}

func quitarExtension(nombreArchivo string) string {
	// Busca la última aparición del punto en el nombre del archivo
	indiceUltimoPunto := strings.LastIndex(nombreArchivo, ".")

	// Si no hay punto, retorna el nombre completo
	if indiceUltimoPunto == -1 {
		return nombreArchivo
	}

	// Retorna el nombre sin la extensión
	return nombreArchivo[:indiceUltimoPunto]
}
func returnParentItemById(id uint64) (parent ParentItemResult) {
	//TODO: Refactorizar este método hay que separar la lógica de la consulta de la lógica de la actualización
	var result ParentItemResult

	item, _ := itemsrepo.NewItemRepository(db.DB).FindByID(id)

	if item.ItemType == "father" {
		result.MainSKU = item.MainSKU
		result.ParentItemID = int(item.ID)

	} else {
		parent, err := itemsparentsrepo.NewItemParentRepository(db.DB).FindByChild(id)
		if err != nil {
			fmt.Printf("Error al ejecutar la consulta: %v", err)
		}
		result.MainSKU = parent.Parent.MainSKU
		result.ParentItemID = int(parent.Parent.ID)

	}

	return result

}

func updateRelatedOrderProcess(fatherOrderId uint64, productEan string, newProductQuantity int64, oldProductQauntity int64, userId uint64, updateType uint64, code string) bool {

	diffProdQuantity := oldProductQauntity - newProductQuantity
	orderItemRepo := orderitemrepo.NewOrderItemRepository(db.DB)
	historicRepo := erpupdateorderlinehistoryrepo.NewErpUpdateOrderLineHistoryRepository(db.DB)
	relationFather, errRelation := suppliersoldorderrelationrepo.NewSupplierSoldOrderRelationRepository(db.DB).FindBySupplierOrder(fatherOrderId)
	//fmt.Println("entra aqui para actualizar")
	if errRelation == gorm.ErrRecordNotFound {
		return false // No se encontró el registro
	}
	// fmt.Println("entra aqui para actualizar 2")
	// fmt.Println("actual quantity: ", oldProductQauntity)
	// fmt.Println("new quantity: ", newProductQuantity)
	// fmt.Println("Diff quantity: ", diffProdQuantity)
	returnedData := returnOrderLinesQuantytiesDataToUpdate(uint64(relationFather.SoldOrderID), productEan)
	// fmt.Println("entra aqui para actualizar 3")
	// fmt.Printf("Datos retornados: %+v\n", returnedData)
	for _, lineData := range returnedData {
		orderItem, _ := orderItemRepo.FindByID(lineData.LineId)
		oldOrderItem := *orderItem
		if lineData.PrepQuantity > 0 {
			if diffProdQuantity <= 0 {
				break
			}

			if lineData.MaxRestQuantity >= int(diffProdQuantity) {
				orderItem.Amount = orderItem.Amount - diffProdQuantity
				diffProdQuantity = 0
			} else {
				orderItem.Amount = orderItem.Amount - int64(lineData.MaxRestQuantity)
				diffProdQuantity = diffProdQuantity - int64(lineData.MaxRestQuantity)
			}
			historicRepo.GenerateOrderLineHistory(oldOrderItem, *orderItem, userId, updateType, code)
			orderItemRepo.Update(orderItem)

		}

	}

	return true
}

type CompareOrderLineQuantyties struct {
	LineId          uint64
	PrepQuantity    int
	PickingQuantity int
	MaxRestQuantity int
}

func returnOrderLinesQuantytiesDataToUpdate(soldFatherId uint64, ean string) []CompareOrderLineQuantyties {

	itemRepo := itemsrepo.NewItemRepository(db.DB)
	orderRepo := orderrepo.NewOrderRepository(db.DB)
	orderItemRepo := orderitemrepo.NewOrderItemRepository(db.DB)

	items, _ := itemRepo.FindByEan(ean)
	orders, _ := orderRepo.FindByFatherId(uint64(soldFatherId))
	itemIds := make([]uint64, 0, len(items))
	orderIds := make([]uint64, 0, len(orders))

	for _, item := range items {
		itemIds = append(itemIds, item.ID)
	}

	for _, order := range orders {
		orderIds = append(orderIds, order.ID)
	}

	pickingMap := make(map[uint64]int) // Suponiendo que ItemID es de tipo uint
	pickingOrders, _ := orderItemRepo.FindByOrderAndItem(orderIds, 1, itemIds, -1, -1)

	for _, pickingOrder := range pickingOrders {
		pickingMap[pickingOrder.ItemID] = int(pickingOrder.Amount)
	}

	preparingOrders, _ := orderItemRepo.FindByOrderAndItem(orderIds, 2, itemIds, -1, -1)
	compareArr := make([]CompareOrderLineQuantyties, 0, len(preparingOrders))

	for _, prepOrder := range preparingOrders {
		compareArr = append(compareArr, CompareOrderLineQuantyties{
			LineId:          prepOrder.ID,
			PrepQuantity:    int(prepOrder.Amount),
			PickingQuantity: pickingMap[prepOrder.ItemID],
			MaxRestQuantity: int(prepOrder.Amount) - pickingMap[prepOrder.ItemID], // Asigna directamente desde el mapa
		})
	}

	return compareArr
}
