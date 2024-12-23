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
