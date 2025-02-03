package erpupdateorderlinehistoryrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type ErpUpdateOrderLineHistoryImpl struct {
	*repositories.GORMRepository[models.ErpUpdateOrderLineHistory]
}

func NewErpUpdateOrderLineHistoryRepository(db *gorm.DB) *ErpUpdateOrderLineHistoryImpl {
	return &ErpUpdateOrderLineHistoryImpl{repositories.NewGORMRepository(db, models.ErpUpdateOrderLineHistory{})}
}

func (repo *ErpUpdateOrderLineHistoryImpl) GenerateOrderLineHistory(orderLine models.OrderItem, ModOrderLine models.OrderItem, userId uint64, updateType uint64, code string) (models.ErpUpdateOrderLineHistory, error) {

	model := models.ErpUpdateOrderLineHistory{
		UpdateErpTypeID:    uint(updateType),
		UpdateGroupCode:    code,
		OrderLineID:        uint(orderLine.ID),
		OrderID:            uint(orderLine.OrderID),
		ItemID:             uint(orderLine.ItemID),
		UserID:             uint(userId),
		StoreID:            uint(orderLine.StoreID),
		Quantity:           int(orderLine.Amount),
		NewQuantity:        int(ModOrderLine.Amount),
		RecivedQuantity:    int(orderLine.RecivedAmount),
		NewRecivedQuantity: int(ModOrderLine.RecivedAmount),
	}
	err := repo.DB.Create(&model).Error

	return model, err
}

type OriginalOrderLine struct {
	OrderLineID uint64
	Quantity    int64
}

func (repo *ErpUpdateOrderLineHistoryImpl) FindByOrders(orders []uint64) ([]OriginalOrderLine, error) {
	var original []OriginalOrderLine
	results := repo.DB.
		Table("erp_update_order_line_histories").
		Select("order_line_id, MAX(quantity) as quantity").
		Where("order_id in ?", orders).
		Group("order_line_id").
		Order("order_line_id DESC").
		Find(&original)

	return original, results.Error
}

type Code struct {
	Code string
}

func (repo *ErpUpdateOrderLineHistoryImpl) FindUpdateCodesByOrders(orders []uint64) ([]Code, error) {
	var codes []Code
	results := repo.DB.
		Table("erp_update_order_line_histories").
		Select("update_group_code as Code").
		Where("order_id in ?", orders).
		Group("update_group_code").
		Order("MIN(created_at) ASC").
		Find(&codes)

	return codes, results.Error
}

func (repo *ErpUpdateOrderLineHistoryImpl) FindHistoryLinesByCode(code string, codes []int) (map[uint64]models.ErpUpdateOrderLineHistory, error) {
	var data []models.ErpUpdateOrderLineHistory
	orderLineMap := make(map[uint64]models.ErpUpdateOrderLineHistory)
	results := repo.DB.
		Where("update_group_code = ?", code).
		Where("update_erp_type_id in ?", codes).
		Find(&data)
	for _, datum := range data {
		orderLineMap[uint64(datum.OrderLineID)] = datum
	}

	return orderLineMap, results.Error
}

type Result struct {
	ModificationDif int    `json:"modification_dif"`
	UserID          uint64 `json:"user_id"`
}

func (repo *ErpUpdateOrderLineHistoryImpl) FindDonePrecentByCode(code string, codes []uint64) ([]Result, error) {

	var results []Result
	err := repo.DB.
		Table("erp_update_order_line_histories").
		Select("SUM(recived_quantity) - SUM(new_recived_quantity) as modification_dif, user_id").
		Where("update_group_code = ?", code).
		Where("update_erp_type_id in ?", codes).
		Group("user_id").
		Find(&results)

	return results, err.Error
}
