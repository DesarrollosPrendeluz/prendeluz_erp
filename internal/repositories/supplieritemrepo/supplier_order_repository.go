package supplieritemrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type SupplierItemRepo interface {
	repositories.Repository[models.SupplierItem]
	FindBySupplierIdAndItemId(itemId uint64, supplierId uint64) (models.SupplierItem, error)
}
