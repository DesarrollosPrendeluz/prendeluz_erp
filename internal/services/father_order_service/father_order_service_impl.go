package services

import (
	"prendeluz/erp/internal/db"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories/fatherorderrepo"
	"prendeluz/erp/internal/repositories/itemsrepo"
	"prendeluz/erp/internal/repositories/orderitemrepo"
)

type FatherOrderImpl struct {
	fatherorderrepo fatherorderrepo.FatherOrderImpl
	itemsRepo       itemsrepo.ItemRepoImpl
	orderitemrepo   orderitemrepo.OrderItemRepoImpl
}

func NewFatherOrderService() *FatherOrderImpl {
	fatherorderrepo := *fatherorderrepo.NewFatherOrderRepository(db.DB)
	itemsRepo := *itemsrepo.NewItemRepository(db.DB)
	orderitemrepo := *orderitemrepo.NewOrderItemRepository(db.DB)

	return &FatherOrderImpl{
		fatherorderrepo: fatherorderrepo,
		itemsRepo:       itemsRepo,
		orderitemrepo:   orderitemrepo}
}

func (s *FatherOrderImpl) FindLinesByFatherOrderCode(pageSize int, offset int, fatherOrderCode string, ean string, supplier_sku string, storeId int) (dtos.FatherOrderOrdersAndLines, int64, error) {
	var result dtos.FatherOrderOrdersAndLines
	var items []models.OrderItem
	var totalRecords int64

	var lines []dtos.LinesInfo
	var itemIds []uint64

	parentData, orderIds, _ := s.fatherorderrepo.FindParentAndOrders(fatherOrderCode)

	itemIds, _ = s.itemsRepo.FindByEanAndSupplierSku(ean, supplier_sku)

	items, totalRecords = s.orderitemrepo.FindByOrderAndItem(orderIds, storeId, itemIds, offset, pageSize)

	//procesado de datos de la query de lineas

	for _, item := range items {
		// Obtener el nombre del proveedor

		supplierName, supplierRef := returnSupplierData(item)
		locations := returnLocations(item)
		var fatherSku string

		if item.Item.ItemType == models.Father {
			fatherSku = item.Item.MainSKU
		} else {
			fatherSku = item.Item.FatherRel.Parent.MainSKU
		}

		// Crear la línea de información
		lineInfo := dtos.LinesInfo{
			LineID:          uint(item.ID),
			OrderCode:       item.OrderID,
			Name:            *item.Item.Name,
			Quantity:        int(item.Amount),
			RecivedQuantity: int(item.RecivedAmount),
			MainSku:         item.Item.MainSKU,
			FatherMainSku:   fatherSku,
			SupplierName:    supplierName,
			SupplierRef:     supplierRef,
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

func returnLocations(item models.OrderItem) []string {
	var locs []string
	var locations *[]models.ItemLocation

	// Determinar qué lista de ubicaciones usar
	if item.Item.ItemType != models.Father && item.Item.FatherRel != nil && item.Item.FatherRel.Parent != nil {
		locations = item.Item.FatherRel.Parent.ItemLocations
	} else {
		locations = item.Item.ItemLocations
	}

	// Recorrer y agregar ubicaciones si existen
	if locations != nil && len(*locations) > 0 {
		for _, location := range *locations {
			locs = append(locs, location.StoreLocations.Name)
		}
	} else {
		locs = append(locs, "")
	}

	return locs
}

func returnSupplierData(item models.OrderItem) (string, string) {
	supplierName, supplierRef := "", ""

	var supplierItems *[]models.SupplierItem

	if item.Item.ItemType != models.Father {
		if item.Item.FatherRel != nil && item.Item.FatherRel.Parent != nil {
			supplierItems = item.Item.FatherRel.Parent.SupplierItems
		}
	} else {
		supplierItems = item.Item.SupplierItems
	}

	if supplierItems != nil && len(*supplierItems) > 0 {
		firstSupplierItem := (*supplierItems)[0]
		if firstSupplierItem.Supplier != nil {
			supplierName = firstSupplierItem.Supplier.Name
			supplierRef = firstSupplierItem.SupplierSku
		}
	}

	return supplierName, supplierRef
}
