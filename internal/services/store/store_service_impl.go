package services

import (
	"log"
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/itemsparentsrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
	"prendeluz/erp/internal/repositories/orderrepo"
	"prendeluz/erp/internal/repositories/storerepo"
	"prendeluz/erp/internal/repositories/storestockrepo"

	"gorm.io/gorm"
)

type StoreServiceImpl struct {
	orderRepo        *orderrepo.OrderRepoImpl
	orderItemsRepo   *orderitemrepo.OrderItemRepoImpl
	itemsRepo        *itemsrepo.ItemRepoImpl
	storeStockRepo   *storestockrepo.StoreStockRepoImpl
	itemsParentsRepo *itemsparentsrepo.ItemsParentsRepoImpl
	storeRepo        *storerepo.StoreRepoImpl
}

func NewStoreService() *StoreServiceImpl {
	orderRepo := orderrepo.NewOrderRepository(db.DB)
	orderItemRepo := orderitemrepo.NewOrderItemRepository(db.DB)
	itemsRepo := itemsrepo.NewItemRepository(db.DB)
	storeStockRepo := storestockrepo.NewStoreStockRepository(db.DB)
	itemsParentsRepo := itemsparentsrepo.NewItemParentRepository(db.DB)
	storeRepo := storerepo.NewStoreRepository(db.DB)

	return &StoreServiceImpl{orderRepo: orderRepo, orderItemsRepo: orderItemRepo, itemsRepo: itemsRepo, storeStockRepo: storeStockRepo, itemsParentsRepo: itemsParentsRepo, storeRepo: storeRepo}
}
func (s *StoreServiceImpl) getParent(child uint64) (models.Item, error) {
	itemsParent, _ := s.itemsParentsRepo.FindByChild(child)
	parent, err := s.itemsRepo.FindByID(itemsParent.ParentItemID)

	return *parent, err
}

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
			itemParent, err := s.getParent(item.ID)
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

func getChilds(items []models.ItemsParents) []models.Item {
	var results []models.Item
	for _, child := range items {
		results = append(results, *child.Child)
	}
	return results
}
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
		childs, _ := s.itemsParentsRepo.FindByParent(itemInStock.ID, 3, 0)
		results = append(results, dtos.ItemStockInfo{Itemname: *itemInStock.Item.Name, SKU: itemInStock.SKU_Parent, Amount: itemInStock.Amount, Childs: getChilds(childs)})
	}

	return results

}
