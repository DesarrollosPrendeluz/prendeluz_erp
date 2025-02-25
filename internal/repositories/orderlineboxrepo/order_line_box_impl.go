package orderlineboxrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type OrderLineBoxImpl struct {
	*repositories.GORMRepository[models.OrderLineBox]
}

func NewOrderLineBoxRepository(db *gorm.DB) *OrderLineBoxImpl {
	return &OrderLineBoxImpl{repositories.NewGORMRepository(db, models.OrderLineBox{})}
}
func (repo *OrderLineBoxImpl) GetByLineId(lineId int) ([]models.OrderLineBox, error) {
	var result []models.OrderLineBox
	err := repo.DB.Where("order_line_id = ?", lineId).Find(&result).Error

	return result, err
}

func (repo *OrderLineBoxImpl) GetByBox(boxId int) ([]models.OrderLineBox, error) {
	var result []models.OrderLineBox
	err := repo.DB.Where("box_id = ?", boxId).Find(&result).Error

	return result, err
}
func (repo *OrderLineBoxImpl) GetOrCreateByOrderLineAndBoxId(orderLineId int, boxId int, quantity int) (models.OrderLineBox, bool, error) {
	var model models.OrderLineBox
	flag := false
	err := repo.DB.
		Where("order_line_id = ?", orderLineId).
		Where("box_id = ?", boxId).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			modelCreate := models.OrderLineBox{
				OrderLineID: orderLineId,
				BoxID:       boxId,
				Quantity:    quantity,
			}
			repo.DB.Create(&modelCreate)
			model = modelCreate
			flag = true

		}
	}
	return model, flag, nil
}
