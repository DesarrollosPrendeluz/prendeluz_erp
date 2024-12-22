package storelocationrepo

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
