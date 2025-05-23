package services

import (
	"fmt"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/repositories/itemsparentsrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/storestockrepo"
)

type ItemsServiceImpl struct {
	itemRepo        itemsrepo.ItemRepoImpl
	orderItemRepo   orderitemrepo.OrderItemRepoImpl
	stockRepo       storestockrepo.StoreStockRepoImpl
	itemParentsRepo itemsparentsrepo.ItemsParentsRepoImpl
}

func NewItemsServiceImpl() *ItemsServiceImpl {
	itemRepo := *itemsrepo.NewItemRepository(db.DB)
	orderItemRepo := *orderitemrepo.NewOrderItemRepository(db.DB)
	stockRepo := *storestockrepo.NewStoreStockRepository(db.DB)
	itemParentsRepo := *itemsparentsrepo.NewItemParentRepository(db.DB)

	return &ItemsServiceImpl{
		itemRepo:        itemRepo,
		orderItemRepo:   orderItemRepo,
		stockRepo:       stockRepo,
		itemParentsRepo: itemParentsRepo,
	}

}

func (s *ItemsServiceImpl) GetItemsForDashboard(flag string, page int, pageSize int) ([]string, int64) {
	var items []string
	var count int64
	if flag == "envio" {
		data, count, _ := s.orderItemRepo.FindByLessOrdered(page, pageSize)
		for _, itemID := range data {
			var sku string
			item, _ := s.itemRepo.FindByID(itemID)

			if item.ItemType != "father" {
				parent, _ := s.itemParentsRepo.FindByChild(item.ID)
				fmt.Println("MELON", parent)
				sku = parent.Parent.MainSKU
			} else {
				sku = item.MainSKU
			}
			items = append(items, sku)
		}
		return items, count

	} else if flag == "stock" {
		items, count, _ = s.stockRepo.FindItemsOrderByQuantity(pageSize, page)
		return items, count

	} else if flag == "coste" {
		data, count, _ := s.itemRepo.FindByPrice(pageSize, page)
		for _, item := range data {
			var sku string
			if item.ItemType != "father" {
				parent, _ := s.itemParentsRepo.FindByChild(item.ID)
				sku = parent.Parent.MainSKU
			} else {
				sku = item.MainSKU
			}
			items = append(items, sku)
		}
		return items, count

	}
	return nil, 0
}
