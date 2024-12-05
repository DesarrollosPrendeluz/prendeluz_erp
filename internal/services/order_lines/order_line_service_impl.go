package services

import (
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
	//var item models.Item
	var label dtos.OrderLineLable
	//itemsRepo := itemsrepo.NewItemRepository(db.DB)
	s.orderItemsRepo.DB.
		Preload("Item").
		Preload("Item.AsinRel.Brand").
		Preload("Item.FatherRel").
		Where("id=?", id).
		First(&orderItem)

	// s.itemsRepo.DB.
	// 	Preload("SupplierItems").
	// 	Preload("SupplierItems.Brand").
	// 	Where("id=?", orderItem.Item.FatherRel.ParentItemID).
	// 	First(&item)

	if orderItem.Item.AsinRel != nil {
		label.Ean = orderItem.Item.AsinRel.Ean
		label.Asin = &orderItem.Item.AsinRel.Code
		label.Brand = orderItem.Item.AsinRel.Brand.Name
		label.BrandAddress = orderItem.Item.AsinRel.Brand.Address
		label.BrandEmail = orderItem.Item.AsinRel.Brand.Email
		label.Company = orderItem.Item.AsinRel.Brand.Company
	}

	// for _, supplierItem := range *item.SupplierItems {

	// 	if supplierItem.Order == 1 {
	// 		label.Brand = supplierItem.Brand.Name
	// 		label.BrandAddress = supplierItem.Brand.Address
	// 		label.BrandEmail = supplierItem.Brand.Email
	// 		label.Company = supplierItem.Brand.Company
	// 		break
	// 	}
	// }

	return label, nil
}
