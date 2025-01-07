package asinrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type AsinRepo interface {
	repositories.Repository[models.Asin]
	FindByItemId(id uint64) (models.Asin, error)
}
