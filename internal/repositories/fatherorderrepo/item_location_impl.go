package fatherorderrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type FatherOrderImpl struct {
	*repositories.GORMRepository[models.FatherOrder]
}

func NewFatherOrderRepository(db *gorm.DB) *FatherOrderImpl {
	return &FatherOrderImpl{repositories.NewGORMRepository(db, models.FatherOrder{})}
}
