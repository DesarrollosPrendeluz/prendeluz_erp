package assignlinerepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type AssignedLineImpl struct {
	*repositories.GORMRepository[models.AssignedLine]
}

func NewAssignedLineImplRepository(db *gorm.DB) *AssignedLineImpl {
	return &AssignedLineImpl{repositories.NewGORMRepository(db, models.AssignedLine{})}
}
