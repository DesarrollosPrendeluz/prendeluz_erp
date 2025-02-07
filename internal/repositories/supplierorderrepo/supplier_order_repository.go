package supplierorderrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type SupplierOrderRepo interface {
	repositories.Repository[models.SupplierOrder]
}
