package ordererrorrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type OrderErrRepo interface {
	repositories.Repository[models.ErrorOrder]
}
