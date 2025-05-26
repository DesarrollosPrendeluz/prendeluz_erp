package orderitemrepo

import (
	"fmt"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type OrderItemRepoImpl struct {
	*repositories.GORMRepository[models.OrderItem]
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepoImpl {
	return &OrderItemRepoImpl{repositories.NewGORMRepository(db, models.OrderItem{})}
}

// Se retornan las lineas de un pedido por el id del pedido
func (repo *OrderItemRepoImpl) FindByOrder(idOrder uint64) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	results := repo.DB.Preload("AssignedRel").Preload("AssignedRel.UserRel").Where("order_id = ?", idOrder).Find(&orderItems)

	return orderItems, results.Error
}

func (repo *OrderItemRepoImpl) FindByOrderAndStore(idOrder uint64, store_id int) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	results := repo.DB.Preload("AssignedRel").Preload("AssignedRel.UserRel").Where("order_id = ?", idOrder).Where("store_id = ?", store_id).Find(&orderItems)

	return orderItems, results.Error
}

// Se retornan las lineas de pedidos las en las cuales el id de los items coincidan
func (repo *OrderItemRepoImpl) FindByItem(idPedido uint64) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	results := repo.DB.Where("item_id = ?", idPedido).Find(&orderItems)

	return orderItems, results.Error
}

// Se retornan las lineas de pedidos las en las cuales el id de los items coincidan
func (repo *OrderItemRepoImpl) FindIdWhereIn(idPedido []uint64) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	results := repo.DB.
		Where("id in ?", idPedido).
		Order("id DESC").
		Find(&orderItems)

	return orderItems, results.Error
}

// // Se retornan las lineas de pedidos las en las cuales el id de los items coincidan
// func (repo *OrderItemRepoImpl) FindWhereIdNotIn(idPedido uint64) ([]models.OrderItem, error) {
// 	var orderItems []models.OrderItem
// 	results := repo.DB.Where("item_id = ?", idPedido).Find(&orderItems)

// 	return orderItems, results.Error
// }

func (repo *OrderItemRepoImpl) FindByOrderExludingIds(ids []uint64, orderId []uint64) ([]models.OrderItem, error) {
	var orderItems []models.OrderItem
	results := repo.DB.Where("id not in ? and order_id in ?", ids, orderId).Find(&orderItems)

	return orderItems, results.Error
}

// Se retornan las lineas de pedidos las en las cuales el id de los items coincidan
func (repo *OrderItemRepoImpl) FindByItemAndOrder(itemId uint64, orderId uint64) (models.OrderItem, error) {
	var orderItems models.OrderItem
	results := repo.DB.Where("item_id = ? and order_id = ?", itemId, orderId).First(&orderItems)

	return orderItems, results.Error
}

func (repo *OrderItemRepoImpl) FindByItemOrderStore(itemId uint64, orderId uint64, storeId uint64) (models.OrderItem, error) {
	var orderItems models.OrderItem
	results := repo.DB.Where("item_id = ? and order_id = ? and store_id = ?", itemId, orderId, storeId).First(&orderItems)

	return orderItems, results.Error
}

// Se retornan las lineas de pedidos las en las cuales el id de los items coincidan
func (repo *OrderItemRepoImpl) FindByItemAndOrders(orderIds []uint64, itemId uint64, storeId uint64) (models.OrderItem, error) {
	var orderItems models.OrderItem
	results := repo.DB.Where("item_id = ? and order_id in ? and store_id = ?", itemId, orderIds, storeId).First(&orderItems)

	return orderItems, results.Error
}

// Se retornan las lineas de pedidos las en las cuales el id de los items coincidan
func (repo *OrderItemRepoImpl) FindByItemsAndOrder(itemIds []uint64, orderId uint64) (models.OrderItem, error) {
	var orderItems models.OrderItem
	results := repo.DB.Where("item_id in ? and order_id = ?", itemIds, orderId).First(&orderItems)

	return orderItems, results.Error
}

func (repo *OrderItemRepoImpl) FindByOrderAndItem(orderIds []uint64, storeId int, itemIds []uint64, offset int, pageSize int) ([]models.OrderItem, int64) {
	var items []models.OrderItem
	var totalRecords int64

	query := addPreloadToShowOrderLineData(repo.DB.Model(&models.OrderItem{}))
	query.
		Where("order_id in ?", orderIds)

	countQuery := repo.DB.
		Model(&models.OrderItem{}).
		Where("order_id in ?", orderIds)

	if storeId != 0 {
		query = query.Where("store_id = ?", storeId)
		countQuery = countQuery.Where("store_id = ?", storeId)
	}

	if len(itemIds) > 0 {
		query.Where("item_id in ?", itemIds)
		countQuery.Where("item_id in ?", itemIds)

	}

	query.Order("CASE WHEN quantity = recived_quantity THEN 1 ELSE 0 END").
		Offset(offset).
		Limit(pageSize).
		Find(&items)

	countQuery.Count(&totalRecords)

	return items, totalRecords
}

func (repo *OrderItemRepoImpl) FindByLineID(lineId []uint64, offset int, pageSize int) ([]models.OrderItem, int64) {
	var items []models.OrderItem
	var totalRecords int64

	query := addPreloadToShowOrderLineData(repo.DB.Model(&models.OrderItem{}))
	query.
		Where("id in ?", lineId).
		Order("CASE WHEN quantity = recived_quantity THEN 1 ELSE 0 END").
		Offset(offset).
		Limit(pageSize).
		Find(&items)

	repo.DB.
		Model(&models.OrderItem{}).
		Where("id in ?", lineId).Count(&totalRecords)

	return items, totalRecords
}

func (repo *OrderItemRepoImpl) FindByLineIDWithOrder(lineId []uint64, order string, offset int, pageSize int) ([]models.OrderItem, int64) {
	var items []models.OrderItem
	var totalRecords int64

	query := addPreloadToShowOrderLineData(repo.DB.Model(&models.OrderItem{}))
	query.
		Where("id in ?", lineId).
		Order("CASE WHEN quantity = recived_quantity THEN 1 ELSE 0 END").
		Order(order).Offset(offset).
		Limit(pageSize).
		Find(&items)

	repo.DB.
		Model(&models.OrderItem{}).
		Where("id in ?", lineId).Count(&totalRecords)

	return items, totalRecords
}
func (repo *OrderItemRepoImpl) FindOrderByIteminPicking(itemId []uint64) models.OrderItem {
	var results models.OrderItem
	repo.DB.Model(&models.OrderItem{}).Where("item_id in ? AND store_id = ?", itemId, 1).
		Order("created_at DESC").
		First(&results)
	return results

}

func (repo *OrderItemRepoImpl) UpdatePickingByItemIdAndOrder(itemId uint64, orderId uint64, quantity int) error {
	results := repo.DB.Model(&models.OrderItem{}).Where("item_id = ? AND order_id = ? AND store_id = 1", itemId, orderId).Update("quantity", quantity)

	return results.Error
}
func (repo *OrderItemRepoImpl) FindByLessOrdered(offset int, pageSize int) ([]uint64, int64, error) {
	type ItemCount struct {
		ItemID uint
		Num    int
	}
	var results []ItemCount
	var count int64
	query := repo.DB.Model(models.OrderItem{}).Distinct("item_id")
	query.Count(&count)
	data := repo.DB.Model(models.OrderItem{}).Select("item_id,count(item_id) as num").Order("num desc").Group("item_id").Offset(offset).Limit(pageSize).Find(&results)
	var ids []uint64
	for _, r := range results {
		ids = append(ids, uint64(r.ItemID))
	}
	return ids, count, data.Error
}
func addPreloadToShowOrderLineData(query *gorm.DB) *gorm.DB {
	return query.
		Preload("Item.FatherRel.Parent.SupplierItems.Supplier").
		Preload("Item.FatherRel.Parent.ItemLocations.StoreLocations").
		Preload("Item.SupplierItems.Supplier").
		Preload("Item.ItemLocations.StoreLocations").
		Preload("AssignedRel.UserRel")
}
