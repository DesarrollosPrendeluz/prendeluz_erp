package supplierrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type SupplierImpl struct {
	*repositories.GORMRepository[models.Supplier]
}

func NewSupplierRepository(db *gorm.DB) *SupplierImpl {
	return &SupplierImpl{repositories.NewGORMRepository(db, models.Supplier{})}
}

func (repo *SupplierImpl) FindAllOrdered(pageSize int, offset int) ([]models.Supplier, error) {
	var storeStocks []models.Supplier

	results := repo.DB.Order("name ASC").Limit(pageSize).Offset(offset).Find(&storeStocks)

	return storeStocks, results.Error

}
