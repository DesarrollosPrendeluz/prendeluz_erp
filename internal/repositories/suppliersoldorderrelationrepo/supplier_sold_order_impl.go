package suppliersoldorderrelationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type SupplierSoldOrderRelationImpl struct {
	*repositories.GORMRepository[models.SupplierSoldOrderRelation]
}

func NewSupplierSoldOrderRelationRepository(db *gorm.DB) *SupplierSoldOrderRelationImpl {
	return &SupplierSoldOrderRelationImpl{repositories.NewGORMRepository(db, models.SupplierSoldOrderRelation{})}
}
