package erpupdateorderlinehistoryrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type ErpUpdateOrderLineHistoryRepo interface {
	repositories.Repository[models.ErpUpdateOrderLineHistory]
}
