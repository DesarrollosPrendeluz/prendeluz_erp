package palletrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type PalletImpl struct {
	*repositories.GORMRepository[models.Pallet]
}

func NewPalletRepository(db *gorm.DB) *PalletImpl {
	return &PalletImpl{repositories.NewGORMRepository(db, models.Pallet{})}
}

func (repo *PalletImpl) GetBoxesAndLinesRaletedDataByOrderId(orderId int, pageSize int, offset int) ([]models.Pallet, error) {
	var models []models.Pallet
	repo.DB.
		Preload("Boxes.BoxContent").
		Where("order_id = ?", orderId).
		Limit(pageSize).
		Offset(offset).
		Find(&models)
	return models, nil
}

func (repo *PalletImpl) GetOrCreatePalletByOrderIdAndNumber(orderId int, number int) (models.Pallet, bool, error) {
	var model models.Pallet
	flag := false
	err := repo.DB.
		Preload("Boxes.BoxContent").
		Where("order_id = ?", orderId).
		Where("number = ?", number).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			modelCreate := models.Pallet{
				OrderID: uint64(orderId),
				Number:  number,
				Label:   "",
			}
			repo.DB.Create(&modelCreate)
			model = modelCreate
			flag = true

		}
	}
	return model, flag, nil
}
