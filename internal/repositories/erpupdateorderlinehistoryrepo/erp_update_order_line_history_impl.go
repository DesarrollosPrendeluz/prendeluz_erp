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
