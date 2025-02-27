package fatherorderrepo

import (
	"fmt"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type FatherOrderImpl struct {
	*repositories.GORMRepository[models.FatherOrder]
}

func NewFatherOrderRepository(db *gorm.DB) *FatherOrderImpl {
	return &FatherOrderImpl{repositories.NewGORMRepository(db, models.FatherOrder{})}
}

func (repo *FatherOrderImpl) FindByCode(code string) (models.FatherOrder, error) {
	var fatherOrder models.FatherOrder
	results := repo.DB.Where("code = ?", code).First(&fatherOrder)

	return fatherOrder, results.Error
}

func (repo *FatherOrderImpl) FindLatestByType(orderType int) (models.FatherOrder, error) {
	var fatherOrder models.FatherOrder
	results := repo.DB.Where("order_type_id = ?", orderType).Order("id DESC").First(&fatherOrder)

	return fatherOrder, results.Error
}

func (repo *FatherOrderImpl) FindAllWithAssocData(pageSize int, offset int, fatherOrderCode string, typeId int, statusId int) ([]dtos.FatherOrderWithRecount, int64, error) {
	var data []dtos.FatherOrderWithRecount
	var results *gorm.DB
	var totalRecords int64

	applyFilters := func(query *gorm.DB, prefix string) *gorm.DB {
		// Filtros de tipo y estado
		if typeId != 0 && statusId != 0 {
			query = query.Where(prefix+"order_type_id = ? AND "+prefix+"order_status_id = ?", typeId, statusId)
		} else if typeId != 0 {
			query = query.Where(prefix+"order_type_id = ?", typeId)
		} else if statusId != 0 {
			query = query.Where(prefix+"order_status_id = ?", statusId)
		}

		// Filtro de código de orden
		if fatherOrderCode != "" {
			query = query.Where(prefix+"code = ?", fatherOrderCode)
		}

		return query
	}

	query := repo.DB.
		Table("father_orders fo").
		Select("fo.id, fo.code, fo.order_status_id, fo.order_type_id, os.name as status, ot.name as type, " +
			"SUM(CASE WHEN ol.store_id = 2 THEN ol.quantity ELSE 0 END) AS total_stock, " +
			"SUM(CASE WHEN ol.store_id = 2 THEN ol.recived_quantity ELSE 0 END) AS pending_stock, " +
			"SUM(CASE WHEN ol.store_id = 1 THEN ol.quantity ELSE 0 END) AS total_picking_stock, " +
			"SUM(CASE WHEN ol.store_id = 1 THEN ol.recived_quantity ELSE 0 END) AS total_recived_picking_quantity").
		Joins("LEFT JOIN orders o ON o.father_order_id = fo.id").
		Joins("LEFT JOIN order_lines ol ON o.id = ol.order_id").
		Joins("LEFT JOIN order_statuses os ON os.id = fo.order_status_id").
		Joins("LEFT JOIN order_types ot ON ot.id = fo.order_type_id")
	query = applyFilters(query, "fo.")
	query = query.Group("fo.id")
	if offset >= 0 && pageSize > 0 {
		query = query.Order("fo.id DESC").Offset(offset).Limit(pageSize)

		results = query.Find(&data)

		query2 := repo.DB.Model(&models.FatherOrder{})
		query2 = applyFilters(query2, "")
		query2.Count(&totalRecords)

		return data, totalRecords, results.Error
	}
	return data, 1, results.Error
}
func (repo *FatherOrderImpl) FindParentAndOrders(code string) (dtos.FatherOrder, []uint64, error) {
	var data models.FatherOrder
	var total uint64
	var partial uint64
	query := repo.DB.
		Preload("SupplierOrder.Supplier").
		Preload("ChildOrders.OrderStatus").
		Preload("OrderStatus").
		Preload("OrderType")

	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	err := query.First(&data).Error

	if err != nil {
		fmt.Println("Error:", err)
	}

	orderIds, orders := func() ([]uint64, []dtos.ChildOrder) {
		var ids []uint64
		var orders []dtos.ChildOrder
		if data.ChildOrders != nil && len(*data.ChildOrders) > 0 {
			for _, order := range *data.ChildOrders {
				ids = append(ids, order.ID)
				orders = append(orders, dtos.ChildOrder{
					ID:              order.ID,
					Code:            order.Code,
					OrderStatusID:   uint(order.OrderStatusID),
					Status:          order.OrderStatus.Name,
					Quantity:        uint64(order.Quantity),
					RecivedQuantity: uint64(order.RecivedQuantity),
				})
				partial = partial + uint64(order.RecivedQuantity)
				total = total + uint64(order.Quantity)

			}

		}
		return ids, orders

	}()
	returnData := dtos.FatherOrder{
		ID:              data.ID,
		Code:            data.Code,
		Type:            data.OrderType.Name,
		OrderTypeID:     uint(data.OrderTypeID),
		Status:          data.OrderStatus.Name,
		OrderStatusID:   uint(data.OrderStatusID),
		Quantity:        total,
		GenericSupplier: &data.SupplierOrder,
		RecivedQuantity: partial,
		Childs:          orders,
	}

	return returnData, orderIds, query.Error
}
