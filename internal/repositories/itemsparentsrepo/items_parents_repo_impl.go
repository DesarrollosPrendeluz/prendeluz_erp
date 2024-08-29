package itemsparentsrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type ItemsParentsRepoImpl struct {
	*repositories.GORMRepository[models.ItemsParents]
}

func NewItemParentRepository(db *gorm.DB) *ItemsParentsRepoImpl {
	return &ItemsParentsRepoImpl{repositories.NewGORMRepository(db, models.ItemsParents{})}
}

func (repo *ItemsParentsRepoImpl) FindByChild(child_id uint64) (models.ItemsParents, error) {
	var item models.ItemsParents
	result := repo.DB.Where("child_item_id = ?", child_id).First(&item)
	return item, result.Error
}
func (repo *ItemsParentsRepoImpl) FindByParent(parent_id uint64, pageSize int, offset int) ([]models.ItemsParents, error) {
	var item []models.ItemsParents
	result := repo.DB.Limit(pageSize).Offset(offset).Preload("Child").Where("parent_item_id = ?", parent_id).Find(&item)
	return item, result.Error
}
