package updateerptyperepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type UpdateErpTypeRepo interface {
	repositories.Repository[models.UpdateErpType]
}
