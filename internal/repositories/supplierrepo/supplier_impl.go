package supplierrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type SupplierImpl struct {
	*repositories.GORMRepository[models.Supplier]
}

func NewSupplierRepository(db *gorm.DB) *SupplierImpl {
	return &SupplierImpl{repositories.NewGORMRepository(db, models.Supplier{})}
}
