package services

import (
	"fmt"
	"io"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"time"

	"prendeluz/erp/internal/repositories/fatherorderrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/ordererrorrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/outorderrelationrepo"
	"prendeluz/erp/internal/utils"
)

type OrderServiceImpl struct {
	orderRepo      orderrepo.OrderRepoImpl
	orderItemsRepo orderitemrepo.OrderItemRepoImpl
	orderErrorRepo repositories.GORMRepository[models.ErrorOrder]
	itemsRepo      itemsrepo.ItemRepoImpl
}

func NewOrderService() *OrderServiceImpl {
	orderRepo := *orderrepo.NewOrderRepository(db.DB)
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})
	orderItemRepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)

	return &OrderServiceImpl{
		orderRepo:      orderRepo,
		orderItemsRepo: orderItemRepo,
		orderErrorRepo: errorOrderRepo,
		itemsRepo:      itemsRepo}
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
			item, exists := itemsMap[orderInfo.MainSku]

			if !exists {
				fmt.Println("error en el sku " + orderInfo.MainSku)
				errorOrder := models.ErrorOrder{
					Main_Sku: orderInfo.MainSku,
					Error:    "Item with sku " + orderInfo.MainSku + " not found",
					Order:    orderCode.OrderCode,
				}

				errorOrdersList = append(errorOrdersList, errorOrder)
			} else {
				orderItem := models.OrderItem{
					ItemID:        item.ID,
					Amount:        orderInfo.Amount,
					RecivedAmount: 0,
					StoreID:       int64(orderInfo.Store),
					ClientID:      orderInfo.Client,
				}
				orderItemsOk[orderCode.OrderCode] = append(orderItemsOk[orderCode.OrderCode], orderItem)
			}
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

// Carga el excel y crea las nuevas ordenes en este caso solo de ventas por el momento
func (s *OrderServiceImpl) UploadOrderExcel(file io.Reader, filename string) error {

	fatherRepo := fatherorderrepo.NewFatherOrderRepository(db.DB)
	fechaActual := time.Now().Format("2006-01-02 15:04:05")

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
		Code:          "OC-" + fechaActual,
		Filename:      "request",
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
func (s *OrderServiceImpl) UploadOrdersByExcel(file io.Reader) error {
	repo := orderrepo.NewOrderRepository(db.DB)
	fatherRepo := fatherorderrepo.NewFatherOrderRepository(db.DB)
	lineRepo := orderitemrepo.NewOrderItemRepository(db.DB)
	itemRepo := itemsrepo.NewItemRepository(db.DB)
	var order models.Order

	excelOrderList, _ := utils.ExceltoJSON(file)
	for _, line := range excelOrderList {

		order, _ = repo.FindByOrderCode(line.OrderCode)
		for _, rowInfo := range line.Info {
			if rowInfo.MainSku != "" {
				item, _ := itemRepo.FindByMainSku(rowInfo.MainSku)
				orderLine, _ := lineRepo.FindByItemAndOrder(item.ID, order.ID)
				orderLine.Amount = rowInfo.Amount
				lineRepo.Update(&orderLine)

			}

		}

		order.OrderStatusID = uint64(orderrepo.Order_Status["en_espera"])
		repo.Update(&order)
		fatherOrder, _ := fatherRepo.FindByID(order.FatherOrderID)
		fatherOrder.OrderStatusID = uint64(orderrepo.Order_Status["en_espera"])
		fatherRepo.Update(fatherOrder)

	}
	return nil

}
