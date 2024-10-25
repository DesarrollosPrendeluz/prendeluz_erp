package outorderrelationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type SupplierImpl struct {
	*repositories.GORMRepository[models.]
}

func NewOutOrderRelationRepository(db *gorm.DB) *OutOrderRelImpl {
	return &OutOrderRelImpl{repositories.NewGORMRepository(db, models.OutOrderRelation{})}
}
