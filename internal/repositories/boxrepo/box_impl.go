package storelocationrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type BoxImpl struct {
	*repositories.GORMRepository[models.Box]
}

func NewBoxRepository(db *gorm.DB) *BoxImpl {
	return &BoxImpl{repositories.NewGORMRepository(db, models.Box{})}
}
