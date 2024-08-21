package itemsrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
)

type ItemRepo interface {
	repositories.Repository[models.Item]
	FindByMainSku(sku string) (models.Item, error)
}

// func NewItemRepository(db *gorm.DB) *ItemRepo {
// 	return &ItemRepo{*NewGORMRepository(db, models.Item{})}
// }
//
// func (repo *ItemRepo) FindByMainSku(sku string) (models.Item, error) {
// 	var item models.Item
//
// 	result := repo.db.Where("main_sku LIKE ?", "%"+sku+"%").First(&item)
//
// 	return item, result.Error
//
// }
