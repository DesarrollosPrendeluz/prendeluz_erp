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

func (repo *FatherOrderImpl) FindAllWithAssocData(pageSize int, offset int, fatherOrderCode string, typeId int, statusId int) ([]dtos.FatherOrderWithRecount, int64, error) {
	var data []dtos.FatherOrderWithRecount
	var results *gorm.DB
	var totalRecords int64

	applyFilters := func(query *gorm.DB) *gorm.DB {
		// Filtros de tipo y estado
		if typeId != 0 && statusId != 0 {
			query = query.Where("fo.order_type_id = ? AND fo.order_status_id = ?", typeId, statusId)
		} else if typeId != 0 {
			query = query.Where("fo.order_type_id = ?", typeId)
		} else if statusId != 0 {
			query = query.Where("fo.order_status_id = ?", statusId)
		}

		// Filtro de código de orden
		if fatherOrderCode != "" {
			query = query.Where("fo.code = ?", fatherOrderCode)
		}

		return query
	}

	query := repo.DB.
		Table("father_orders fo").
		Select("fo.id, fo.code, fo.order_status_id, os.name as status, ot.name as type, fo.order_type_id, SUM(ol.quantity) as total_stock, SUM(ol.recived_quantity) as pending_stock").
		Joins("LEFT JOIN orders o ON o.father_order_id = fo.id").
		Joins("LEFT JOIN order_lines ol ON o.id = ol.order_id").
		Joins("LEFT JOIN order_statuses os ON os.id = fo.order_status_id").
		Joins("LEFT JOIN order_types ot ON ot.id = fo.order_type_id")
	query = applyFilters(query)
	query = query.Group("fo.id")
	if offset >= 0 && pageSize > 0 {
		query = query.Offset(offset).Limit(pageSize)
	}
	results = query.Find(&data)

	query2 := repo.DB.Model(&models.FatherOrder{})
	query2 = applyFilters(query2)
	query2.Count(&totalRecords)

	return data, totalRecords, results.Error
}

func (repo *FatherOrderImpl) FindLinesByFatherOrderCode(pageSize int, offset int, fatherOrderCode string, ean string) (dtos.FatherOrderOrdersAndLines, int64, error) {
	var result dtos.FatherOrderOrdersAndLines
	var items []models.OrderItem
	var totalRecords int64
	var results []models.Item
	var lines []dtos.LinesInfo

	parentData, orderIds, _ := repo.findParentAndOrders(fatherOrderCode)

	//query de obtención de datos de lineas

	query := repo.DB.
		Model(&models.OrderItem{}).
		Preload("AssignedRel.UserRel").
		Preload("Item.FatherRel.Parent.SupplierItems.Supplier").
		Preload("Item.FatherRel.Parent.ItemLocations.StoreLocations").
		Where("order_id in ?", orderIds)
	countQuery := repo.DB.
		Model(&models.OrderItem{}).
		Where("order_id in ?", orderIds)
	if ean != "" {
		errr := repo.DB.
			Table("items AS i").
			Where("i.ean = ?", ean).
			Find(&results).Error

		itemIds := func() []uint64 {
			var ids []uint64
			for _, order := range results {
				ids = append(ids, order.ID)

			}
			return ids

		}()
		if errr != nil {
			fmt.Println(errr.Error())

		}
		query.Where("item_id in ?", itemIds)
		countQuery.Where("item_id in ?", itemIds)

	}

	query.
		Offset(offset).
		Limit(pageSize).
		Find(&items)

	countQuery.
		Count(&totalRecords)

	//procesado de datos de la query de lineas

	for _, item := range items {
		// Obtener el nombre del proveedor
		supplierName := func() string {
			if item.Item.FatherRel.Parent.SupplierItems != nil && len(*item.Item.FatherRel.Parent.SupplierItems) > 0 {
				return (*item.Item.FatherRel.Parent.SupplierItems)[0].Supplier.Name
			}
			return ""
		}()

		// Obtener las ubicaciones
		locations := func() []string {
			var locs []string
			if item.Item.FatherRel.Parent.ItemLocations != nil && len(*item.Item.FatherRel.Parent.ItemLocations) > 0 {
				for _, location := range *item.Item.FatherRel.Parent.ItemLocations {
					locs = append(locs, location.StoreLocations.Name)
				}
			}
			return locs
		}()

		// Crear la línea de información
		lineInfo := dtos.LinesInfo{
			LineID:          uint(item.ID),
			OrderCode:       item.OrderID,
			Name:            *item.Item.Name,
			Quantity:        int(item.Amount),
			RecivedQuantity: int(item.RecivedAmount),
			MainSku:         item.Item.MainSKU,
			Ean:             item.Item.EAN,
			SupplierName:    supplierName,
			Location:        locations,
			AssignedUser: dtos.AssignedUserToOrderItem{
				AssignationId: item.AssignedRel.ID,
				UserId:        item.AssignedRel.UserID,
				UserName:      item.AssignedRel.UserRel.Name,
			},
		}

		// Añadir la línea al resultado
		lines = append(lines, lineInfo)
	}
	//Monatje de lineas
	result.FatherOrder = parentData
	result.Lines = lines

	return result, totalRecords, nil
}

func (repo *FatherOrderImpl) findParentAndOrders(code string) (dtos.FatherOrder, []uint64, error) {
	var data models.FatherOrder
	var total uint64
	var partial uint64
	query := repo.DB.
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
		RecivedQuantity: partial,
		Childs:          orders,
	}

	return returnData, orderIds, query.Error
}
