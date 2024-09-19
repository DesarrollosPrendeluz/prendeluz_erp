package inorderrelationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type InOrderRelImpl struct {
	*repositories.GORMRepository[models.InOrderRelation]
}

func NewInOrderRelationRepository(db *gorm.DB) *InOrderRelImpl {
	return &InOrderRelImpl{repositories.NewGORMRepository(db, models.InOrderRelation{})}
}
