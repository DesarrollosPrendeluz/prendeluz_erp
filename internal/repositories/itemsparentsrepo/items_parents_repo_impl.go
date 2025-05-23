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

// Busca un producto hijo en base a su aparición en la tabla parent_items
func (repo *ItemsParentsRepoImpl) FindByChild(child_id uint64) (models.ItemsParents, error) {
	var item models.ItemsParents
	result := repo.DB.Preload("Parent").Where("child_item_id = ?", child_id).First(&item)
	return item, result.Error
}

// Busca un producto padre en base a su aparición en la tabla parent_items y precarga los hijos
func (repo *ItemsParentsRepoImpl) FindByParent(parent_id uint64, pageSize int, offset int) ([]models.ItemsParents, error) {
	var item []models.ItemsParents
	result := repo.DB.Limit(pageSize).Offset(offset).Preload("Child.AsinRel").Where("parent_item_id = ?", parent_id).Find(&item)
	return item, result.Error
}

func (repo *ItemsParentsRepoImpl) FindMultipleParents(childs []uint64) ([]uint64, error) {
	var parents []uint64
	result := repo.DB.Preload("items").Select("parent_item_id").Where("child_item_id in ?", childs).Find(&parents)
	return parents, result.Error
}
