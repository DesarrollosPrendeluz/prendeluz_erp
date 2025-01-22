package services

import "prendeluz/erp/internal/models"

type StadisitcsService interface {
	GetChangeStadistics(fatherId uint64) []models.OrderItem
}
