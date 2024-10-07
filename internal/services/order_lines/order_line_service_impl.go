package services

import (
	"fmt"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
)

type OrderLineServiceImpl struct {
	orderRepo      orderrepo.OrderRepoImpl
	orderItemsRepo orderitemrepo.OrderItemRepoImpl
	//itemRepo       itemsrepo.ItemRepoImpl
	orderErrorRepo repositories.GORMRepository[models.ErrorOrder]
	itemsRepo      itemsrepo.ItemRepoImpl
}

func NewOrderLineServiceImpl() *OrderLineServiceImpl {
	orderRepo := *orderrepo.NewOrderRepository(db.DB)
	errorOrderRepo := *repositories.NewGORMRepository(db.DB, models.ErrorOrder{})
	orderItemRepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)

	return &OrderLineServiceImpl{
		orderRepo:      orderRepo,
		orderItemsRepo: orderItemRepo,
		orderErrorRepo: errorOrderRepo,
		itemsRepo:      itemsRepo,
	}
}

func (s *OrderLineServiceImpl) OrderLineLabel(id int) (dtos.OrderLineLable, error) {
	var orderItem models.OrderItem
	var item models.Item
	var label dtos.OrderLineLable
	//itemsRepo := itemsrepo.NewItemRepository(db.DB)
	s.orderItemsRepo.DB.Debug().
		Preload("Item").
		Preload("Item.AsinRel").
		Preload("Item.FatherRel").
		Where("id=?", id).
		Find(&orderItem)
	fmt.Println(orderItem.Item.FatherRel.ParentItemID)
	fmt.Println("Query dos en proceso")

	s.itemsRepo.DB.Debug().
		Preload("SupplierItems").
		Preload("SupplierItems.Brand").
		Where("id=?", orderItem.Item.FatherRel.ParentItemID).
		Find(&item)

	label.Ean = orderItem.Item.EAN
	label.Asin = orderItem.Item.AsinRel.Code

	for _, supplierItem := range *item.SupplierItems {

		if supplierItem.Order == 1 {
			label.Brand = supplierItem.Brand.Name
			label.BrandAddress = supplierItem.Brand.Address
			label.BrandEmail = supplierItem.Brand.Email
			break
		}
	}

	return label, nil
}
