package services

import "prendeluz/erp/internal/models"

type ErpUpdateTypesService interface {
	GetAll() []models.UpdateErpType
}
