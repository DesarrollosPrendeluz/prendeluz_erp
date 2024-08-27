package services

import (
	"io"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/utils"

	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	orderRepo      repositories.GORMRepository[models.Order]
	orderItemsRepo orderitemrepo.OrderItemRepoImpl
	orderErrorRepo repositories.GORMRepository[models.ErrorOrder]
	itemsRepo      itemsrepo.ItemRepoImpl
}

func NewOrderService() *OrderServiceImpl {
	orderRepo := *repositories.NewGORMRepository(db.DB, models.Order{})
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})
	orderItemRepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)

	return &OrderServiceImpl{orderRepo: orderRepo, orderItemsRepo: orderItemRepo,
		orderErrorRepo: errorOrderRepo, itemsRepo: itemsRepo}
}

func itemsExist(rawOrders []utils.ExcelOrder, filename string) ([]models.Order, map[string][]models.OrderItem, []models.ErrorOrder) {
	itemsRepo := itemsrepo.NewItemRepository(db.DB)

	var errorOrdersList []models.ErrorOrder
	var ordersList []models.Order
	orderItemsOk := make(map[string][]models.OrderItem)

	for _, orderCode := range rawOrders {
		for _, orderInfo := range orderCode.Info {
			item, err := itemsRepo.FindByMainSku(orderInfo.MainSku)

			if err != nil {
				errorOrder := models.ErrorOrder{
					Main_Sku: orderInfo.MainSku,
					Error:    "Item with sku not found",
					Order:    orderCode.OrderCode,
				}

				errorOrdersList = append(errorOrdersList, errorOrder)
			} else {
				orderItem := models.OrderItem{
					ItemID: item.ID,
					Amount: orderInfo.Amount,
				}
				orderItemsOk[orderCode.OrderCode] = append(orderItemsOk[orderCode.OrderCode], orderItem)
			}
		}

		order := models.Order{
			Orden_compra: orderCode.OrderCode,
			Filename:     filename,
		}
		ordersList = append(ordersList, order)

	}
	return ordersList, orderItemsOk, errorOrdersList
}
func (s *OrderServiceImpl) UploadOrderExcel(file io.Reader, filename string) error {

	excelOrderList, err := utils.ExceltoJSON(file)
	if err != nil {
		return err
	}
	succesOrders, orderItems, errorOrders := itemsExist(excelOrderList, filename)

	return db.DB.Transaction(func(tx *gorm.DB) error {
		s.orderRepo.SetDB(tx)
		s.orderItemsRepo.SetDB(tx)
		s.orderErrorRepo.SetDB(tx)
		_, err := s.orderRepo.CreateAll(&succesOrders)
		if err != nil {
			return err
		}

		if len(errorOrders) > 0 {
			_, err := s.orderErrorRepo.CreateAll(&errorOrders)
			if err != nil {
				return err
			}
		}
		var orderItemToInsert []models.OrderItem
		for _, order := range succesOrders {
			for _, tmp := range orderItems[order.Orden_compra] {
				tmp.OrderID = order.ID
				orderItemToInsert = append(orderItemToInsert, tmp)
			}
		}
		_, err = s.orderItemsRepo.CreateAll(&orderItemToInsert)
		if err != nil {
			return err
		}
		return nil

	})
}

func (s *OrderServiceImpl) GetOrders() ([]dtos.ItemsPerOrder, error) {

	var results []dtos.ItemsPerOrder

	orders, err := s.orderRepo.FindAll()
	if err != nil {
		return results, err
	}

	for _, order := range orders {
		var itemOrder dtos.ItemsPerOrder
		itemOrder.OrderCode = order.Orden_compra
		orderItemList, _ := s.orderItemsRepo.FindByOrder(order.ID)

		for _, orderItem := range orderItemList {
			var itemInfo dtos.ItemInfo
			item, _ := s.itemsRepo.FindByID(orderItem.ItemID)
			itemInfo.Amount = orderItem.Amount
			itemInfo.Sku = item.MainSKU
			itemOrder.ItemsOrdered = append(itemOrder.ItemsOrdered, itemInfo)
		}
		results = append(results, itemOrder)

	}

	return results, nil

}
