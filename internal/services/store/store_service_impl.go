package services

import (
	"log"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/storestockrepo"

	"gorm.io/gorm"
)

type StoreServiceImpl struct {
	orderRepo      *repositories.GORMRepository[models.Order]
	orderItemsRepo *orderitemrepo.OrderItemRepoImpl
	itemsRepo      *itemsrepo.ItemRepoImpl
	storeStockRepo *storestockrepo.StoreStockRepoImpl
}

func NewStoreService() *StoreServiceImpl {
	orderRepo := repositories.NewGORMRepository(db.DB, models.Order{})
	orderItemRepo := orderitemrepo.NewOrderItemRepository(db.DB)
	itemsRepo := itemsrepo.NewItemRepository(db.DB)
	storeStockRepo := storestockrepo.NewStoreStockRepository(db.DB)

	return &StoreServiceImpl{orderRepo: orderRepo, orderItemsRepo: orderItemRepo, itemsRepo: itemsRepo, storeStockRepo: storeStockRepo}
}

func (s *StoreServiceImpl) UpdateStoreStock(order_id uint64) error {
	itemsOrdered := make(map[string]int64)
	orders, _ := s.orderItemsRepo.FindByOrder(10)
	for _, order := range orders {
		item, _ := s.itemsRepo.FindByID(order.ItemID)
		if item.ItemType != "father" {
			log.Println("Skip")
		} else {
			itemsOrdered[item.MainSKU] += order.Amount
		}
	}
	return db.DB.Transaction(func(tx *gorm.DB) error {
		s.itemsRepo.SetDB(tx)
		s.orderItemsRepo.SetDB(tx)
		s.orderRepo.SetDB(tx)
		s.storeStockRepo.SetDB(tx)

		var updateStock []models.StoreStock

		for item, amount := range itemsOrdered {
			itemToUpdate, _ := s.storeStockRepo.FindByItem(item)

			itemToUpdate.Amount -= amount
			updateStock = append(updateStock, itemToUpdate)
		}
		log.Println(updateStock)
		s.storeStockRepo.UpdateAll(&updateStock)
		return nil
	})
}
