package stockdeficitrepo

import (
	"errors"
	"fmt"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type StockDeficitImpl struct {
	*repositories.GORMRepository[models.StockDeficit]
}
type ParentItemResult struct {
	ParentItemID int    `gorm:"column:parent_item_id"`
	MainSKU      string `gorm:"column:main_sku"`
}
type StockDeficitResult struct {
	Deficit float64 `gorm:"column:deficit"`
}

func NewStockDeficitRepository(db *gorm.DB) *StockDeficitImpl {
	return &StockDeficitImpl{repositories.NewGORMRepository(db, models.StockDeficit{})}
}

func (repo *StockDeficitImpl) GetallByStore(storeId int, pageSize int, offset int) ([]models.StockDeficit, error) {
	var models []models.StockDeficit
	repo.DB.
		Preload("Item.SupplierItems.Supplier").
		Where("store_id = ?", storeId).
		Where("quantity != 0 OR pending_stock != 0").
		Limit(pageSize).
		Offset(offset).
		Find(&models)
	return models, nil
}
func (repo *StockDeficitImpl) GetallByStoreAndSupplier(storeId int, supplier int, pageSize int, offset int) ([]models.StockDeficit, error) {
	var modelsData []models.StockDeficit

	subQuery := repo.DB.
		Model(&models.StockDeficit{}).
		Select("stock_deficits.id").
		Joins("JOIN items ON items.main_sku = stock_deficits.parent_main_sku").
		Joins("JOIN supplier_items ON supplier_items.item_id = items.id").
		//Where("supplier_items.item_id = ?", 2).
		Where("store_id = ?", storeId).
		Where("quantity != 0 OR pending_stock != 0").
		Where("supplier_items.supplier_id = ?", supplier)

	err := repo.DB.
		//("Item.SupplierItems", "item_id = ?", 2).
		Preload("Item.SupplierItems.Supplier").
		Where("id IN (?)", subQuery).
		Limit(pageSize).
		Offset(offset).
		Find(&modelsData).Error

	if err != nil {
		return nil, err
	}
	return modelsData, nil
}

func (repo *StockDeficitImpl) CountConditional(storeId int) (int64, error) {
	var count int64
	err := repo.DB.Table("stock_deficits").Count(&count).Where("store_id = ?", storeId).Error
	return count, err
}

func (repo *StockDeficitImpl) CalcStockDeficitByItem(child_item_id uint64, store_id int64) {
	var result ParentItemResult
	var existing models.StockDeficit
	var deficit StockDeficitResult
	var pending StockDeficitResult

	if err := repo.DB.Table("item_parents AS ip").
		Select("ip.parent_item_id, i.main_sku").
		Joins("JOIN items i ON i.id = ip.parent_item_id").
		Where("ip.child_item_id = ?", child_item_id).
		Limit(1).
		Scan(&result).Error; err != nil {
		fmt.Printf("Error al ejecutar la consulta: %v", err)
	}

	err2 := repo.DB.
		Table("order_lines AS ol").
		Select("GREATEST(0, -(IFNULL(AVG(ss.quantity), 0) - (SUM(ol.quantity) - SUM(ol.recived_quantity)))) AS deficit").
		Joins("LEFT JOIN item_parents ip ON ip.child_item_id = ol.item_id").
		Joins("LEFT JOIN orders ord ON ord.id = ol.order_id").
		Joins("INNER JOIN father_orders fo ON fo.id = ord.father_order_id AND fo.order_type_id = 2 AND fo.order_status_id != 3").
		Joins("LEFT JOIN store_stocks ss ON ss.parent_main_sku = ?", result.MainSKU).
		Where("ip.parent_item_id = ?", result.ParentItemID).
		Where("ol.store_id = ?", store_id).
		Group("ip.parent_item_id").
		Take(&deficit).Error
	if err2 != nil {
		fmt.Printf("Error al ejecutar la consulta: %v", err2)
	}
	err2 = repo.DB.
		Table("order_lines AS ol").
		Select(" SUM(ol.quantity) - SUM(ol.recived_quantity) AS deficit").
		Joins("LEFT JOIN item_parents ip ON ip.child_item_id = ol.item_id").
		Joins("LEFT JOIN orders ord ON ord.id = ol.order_id").
		Joins("INNER JOIN father_orders fo ON fo.id = ord.father_order_id AND fo.order_type_id = 1 AND fo.order_status_id != 3").
		Where("ip.parent_item_id = ?", result.ParentItemID).
		Where("ol.store_id = ?", store_id).
		Group("ip.parent_item_id").
		Take(&pending).Error

	if err2 != nil {
		fmt.Printf("Error al ejecutar la consulta: %v", err2)
	}

	err3 := repo.DB.Table("stock_deficits").
		Where("store_id = ? AND parent_main_sku = ?", store_id, result.MainSKU).
		First(&existing).Error

	if err3 != nil {
		if errors.Is(err3, gorm.ErrRecordNotFound) {
			// El registro no existe, realizar una inserción
			newRecord := models.StockDeficit{
				StoreID:       uint64(store_id),
				SKU_Parent:    result.MainSKU,
				Amount:        int64(deficit.Deficit),
				PendingAmount: int64(pending.Deficit),
			}
			repo.Create(&newRecord)

		} else {
			// Error diferente a "registro no encontrado"
			fmt.Errorf("error al buscar registro existente: %w", err3)
		}
	} else {
		// El registro ya existe, realizar una actualización
		existing.Amount = int64(deficit.Deficit)
		existing.PendingAmount = int64(pending.Deficit)
		repo.Update(&existing)

	}

}
