package supplierrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type SupplierRepo interface {
	repositories.Repository[models.Supplier]
}
