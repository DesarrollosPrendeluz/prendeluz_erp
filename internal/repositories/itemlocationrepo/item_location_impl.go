package itemlocationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type ItemLocationImpl struct {
	*repositories.GORMRepository[models.ItemLocation]
}

func NewInItemLocationRepository(db *gorm.DB) *ItemLocationImpl {
	return &ItemLocationImpl{repositories.NewGORMRepository(db, models.ItemLocation{})}
}
