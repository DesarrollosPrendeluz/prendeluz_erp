package supplierorderrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type SupplierOrderImpl struct {
	*repositories.GORMRepository[models.SupplierOrder]
}

func NewSupplierOrderRepository(db *gorm.DB) *SupplierOrderImpl {
	return &SupplierOrderImpl{repositories.NewGORMRepository(db, models.SupplierOrder{})}
}
