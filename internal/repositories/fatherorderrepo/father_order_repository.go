package fatherorderrepo

import (
	"prendeluz/erp/internal/dtos"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type FatherOrderRepo interface {
	repositories.Repository[models.FatherOrder]
	FindAllWithAssocData(pageSize int, offset int) ([]dtos.FatherOrderWithRecount, error)
	FindParentAndOrders(code string) (dtos.FatherOrder, []uint64, error)
	FindByCode(code string) (models.FatherOrder, error)
}
