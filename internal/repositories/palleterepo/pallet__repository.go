package storelocationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type PalletRepo interface {
	repositories.Repository[models.Pallet]
}
