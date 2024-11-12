package fatherorderrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type FatherOrderRepo interface {
	repositories.Repository[models.FatherOrder]
}
