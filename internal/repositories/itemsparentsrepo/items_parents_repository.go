package itemsparentsrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type ItemsParentsRepo interface {
	repositories.Repository[models.Item]
	FindByChild(child_id uint64) (models.ItemsParents, error)
	FindByParent(parent_id uint64, pageSize int, offset int) ([]models.ItemsParents, error)
}
