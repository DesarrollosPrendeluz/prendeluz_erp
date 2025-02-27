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

func (repo *SupplierSoldOrderRelationImpl) FindBySupplierOrder(id uint64) (models.SupplierSoldOrderRelation, error) {
	var relation models.SupplierSoldOrderRelation

	result := repo.DB.Debug().Where("supplier_father_order_id = ?", id).First(&relation)

	return relation, result.Error

}
