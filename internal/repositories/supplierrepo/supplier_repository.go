package supplierrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type SupplierRepo interface {
	repositories.Repository[models.Supplier]
	FindAllOrdered(pageSize int, offset int) ([]models.Supplier, error)
}
