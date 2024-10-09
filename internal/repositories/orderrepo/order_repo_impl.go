package orderrepo

import (
	"fmt"
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

var Order_Status = map[string]int{
	"iniciada":   1,
	"finalizada": 2,
	"en_proceso": 3,
	"en_espera":  4,
}

var Order_Types = map[string]int{
	"compra": 1,
	"venta":  2,
}

type OrderRepoImpl struct {
	*repositories.GORMRepository[models.Order]
}

func NewOrderRepository(db *gorm.DB) *OrderRepoImpl {
	return &OrderRepoImpl{repositories.NewGORMRepository(db, models.Order{})}

}

// Busca una orden por el codigo alfanumerico asociado
func (repo *OrderRepoImpl) FindByOrderCode(orderCode string) (models.Order, error) {
	var order models.Order

	results := repo.DB.Where("code LIKE ?", "%"+orderCode+"%").First(&order)

	return order, results.Error
}

// FindOrderByDate recupera un pedido basado en un rango de fechas.
// Acepta startDate y endDate como cadenas de texto con el formato: "YYYY-MM-DD".
// Si se proporcionan tanto startDate como endDate, devuelve los pedidos entre esas fechas.
// Si solo se proporciona startDate, devuelve los pedidos para esa fecha específica.
// Si solo se proporciona endDate, devuelve los pedidos para esa fecha específica.
// Devuelve una estructura Order y un error si algo sale mal.

func (repo *OrderRepoImpl) FindOrderFiltered(pageSize int, page int, startDate string, endDate string, typeId int, statusId int, orderCode string) ([]models.Order, int64, error) {
	var orders []models.Order
	var totalRecords int64
	var results *gorm.DB

	applyFilters := func(query *gorm.DB) *gorm.DB {
		// Filtros de fecha
		if startDate != "" && endDate != "" {
			query = query.Where("date(created_at) BETWEEN ? AND ?", startDate, endDate)
		} else if startDate != "" {
			query = query.Where("date(created_at) = ?", startDate)
		} else if endDate != "" {
			query = query.Where("date(created_at) = ?", endDate)
		}

		// Filtros de tipo y estado
		if typeId != 0 && statusId != 0 {
			query = query.Where("order_type_id = ? AND order_status_id = ?", typeId, statusId)
		} else if typeId != 0 {
			query = query.Where("order_type_id = ?", typeId)
		} else if statusId != 0 {
			query = query.Where("order_status_id = ?", statusId)
		}

		// Filtro de código de orden
		if orderCode != "" {
			query = query.Where("code = ?", orderCode)
		}

		return query
	}
	query := repo.DB.Debug().Preload("OrderStatus").Preload("OrderType")
	query = applyFilters(query)

	query2 := repo.DB.Model(&models.Order{})
	query2 = applyFilters(query2)

	// Obtener el total de registros sin paginación
	query2.Count(&totalRecords)
	//totalRecords = 1
	// Agregar paginación
	if page >= 0 && pageSize > 0 {
		query = query.Offset(page).Limit(pageSize)
	}

	// Ejecutar la consulta paginada
	fmt.Println(page)
	fmt.Println(pageSize)
	results = query.Find(&orders)
	fmt.Printf("%+v\n", orders)
	fmt.Println(totalRecords)

	return orders, totalRecords, results.Error
}

func (r *OrderRepoImpl) FindAll(pageSize int, offset int) ([]models.Order, int64, error) {
	var items []models.Order
	var totalRecords int64

	// Primero obtener el recuento total de registros
	result := r.DB.Model(&models.Order{}).Count(&totalRecords)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Luego obtener los registros paginados
	result = r.DB.Preload("OrderStatus").Preload("OrderType").Limit(pageSize).Offset(offset).Find(&items)

	return items, totalRecords, result.Error
}

// Actualiza el estado de una orden
// Recibe el id del nuevo estado y el id de la orden
func (repo *OrderRepoImpl) UpdateStatus(newStatus int, orderID uint64) error {
	results := repo.DB.Model(models.Order{}).Where("id = ?", orderID).Update("status_id", newStatus)

	return results.Error

}

func (repo *OrderRepoImpl) GetSupplierOrders(order_type *int) ([]dtos.SupplierOrders, error) {
	var orders []dtos.SupplierOrders

	// Consulta SQL manual con JOIN
	query := `
		SELECT 
			o.id as order_code, 
			orl.quantity as stock_to_buy, 
			it.main_sku as item_sku, 
			it.id as item_id,
			ip.parent_item_id as father_id,
			it.name as name,
			it.ean as ean,
			sp.name as supplier_name,
			spi.supplier_sku as supplier_code,
			spi.price as supplier_price
		FROM orders  o
		INNER JOIN order_lines as orl ON orl.order_id = o.id 
		LEFT JOIN items as it ON it.id = orl.item_id
		LEFT JOIN item_parents ip on ip.child_item_id = it.id
		LEFT JOIN supplier_items as spi ON spi.item_id = ip.parent_item_id AND spi.order = 1
		LEFT JOIN suppliers as sp ON sp.id = spi.supplier_id
		WHERE o.order_type_id = 2
		
	`
	if order_type != nil {
		query += " AND o.order_type_id = ?"
	}

	query += " ORDER BY o.id"

	// Ejecutamos la consulta con Raw y mapeamos los resultados al slice de `orders`
	if err := repo.DB.Raw(query).Scan(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}
