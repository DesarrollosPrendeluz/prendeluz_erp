package supplieritemrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type SupplierItemImpl struct {
	*repositories.GORMRepository[models.SupplierItem]
}

func NewSupplierItemRepository(db *gorm.DB) *SupplierItemImpl {
	return &SupplierItemImpl{repositories.NewGORMRepository(db, models.SupplierItem{})}
}

func (repo *SupplierItemImpl) FindBySupplierIdAndItemId(itemId uint64, supplierId uint64) (models.SupplierItem, error) {
	var item models.SupplierItem

	result := repo.DB.Where("item_id = ?", itemId).Where("supplier_id = ?", supplierId).First(&item)

	return item, result.Error

}
