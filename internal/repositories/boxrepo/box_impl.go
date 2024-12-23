package boxrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type BoxImpl struct {
	*repositories.GORMRepository[models.Box]
}

func NewBoxRepository(db *gorm.DB) *BoxImpl {
	return &BoxImpl{repositories.NewGORMRepository(db, models.Box{})}
}

func (repo *BoxImpl) GetOrCreateBoxByPalletIdAndNumber(palletId int, number int, quantity int) (models.Box, bool, error) {
	var model models.Box
	flag := false
	err := repo.DB.
		Where("pallet_id = ?", palletId).
		Where("number = ?", number).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			modelCreate := models.Box{
				PalletID: uint64(palletId),
				Number:   number,
				Label:    "",
				Quantity: quantity,
			}
			repo.DB.Create(&modelCreate)
			model = modelCreate
			flag = true

		}
	}
	return model, flag, nil
}
