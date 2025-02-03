package services

import (
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
)

type PalletService interface {
	Get(pallet int, page int, pageSize int) ([]models.Pallet, int64, error)
	GetPalletByOrder(orderId int, page int, pageSize int) ([]models.Pallet, int64, error)
	Create(data dtos.PalletCreateReq) []error
	Update(data dtos.PalletUpdateReq) []error
}
