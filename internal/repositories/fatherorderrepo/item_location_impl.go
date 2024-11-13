package fatherorderrepo

import (
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

	query := repo.DB.Debug().
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
