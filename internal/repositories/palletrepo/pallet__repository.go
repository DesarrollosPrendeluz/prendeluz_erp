package palletrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type PalletRepo interface {
	repositories.Repository[models.Pallet]
	GetBoxesAndLinesRaletedDataByOrderId(orderId int, pageSize int, offset int) ([]models.Pallet, error)
}
