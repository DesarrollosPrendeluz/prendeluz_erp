package suppliersoldorderrelationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type SupplierSoldOrderRepo interface {
	repositories.Repository[models.SupplierSoldOrderRelation]
}
