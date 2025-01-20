package updateerptyperepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type UpdateErpTypeImpl struct {
	*repositories.GORMRepository[models.UpdateErpType]
}

func NewUpdateErpTypeRepository(db *gorm.DB) *UpdateErpTypeImpl {
	return &UpdateErpTypeImpl{repositories.NewGORMRepository(db, models.UpdateErpType{})}
}
