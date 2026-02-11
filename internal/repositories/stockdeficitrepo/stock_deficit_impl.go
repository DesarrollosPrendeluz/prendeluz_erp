package stockdeficitrepo

import (
	"fmt"
	"log"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type StockDeficitImpl struct {
	*repositories.GORMRepository[models.StockDeficit]
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

func (repo *StockDeficitImpl) GetByFatherAndStore(fatherSku string, store int64) (models.StockDeficit, error) {
	var modelsData models.StockDeficit

	err := repo.DB.
		Where("parent_main_sku = ?", fatherSku).
		Where("store_id = ?", store).
		First(&modelsData).Error

	return modelsData, err

}

func (repo *StockDeficitImpl) FindOrCreateByFatherAndStore(fatherSku string, store int64) (models.StockDeficit, error) {
	var modelsData models.StockDeficit

	err := repo.DB.
		Where("parent_main_sku = ?", fatherSku).
		Where("store_id = ?", store).
		First(&modelsData)
	if err != nil {
		if err.Error == gorm.ErrRecordNotFound {
			modelCreate := models.StockDeficit{
				SKU_Parent:    fatherSku,
				StoreID:       uint64(store),
				Amount:        0,
				PendingAmount: 0,
			}
			repo.DB.Create(&modelCreate)
			modelsData = modelCreate

		}
	}

	return modelsData, err.Error

}
func (repo *StockDeficitImpl) GetByRegsitersByFatherSkuIn(filter []string, store int, page int, pageSize int) ([]models.StockDeficit, error) {
	var modelsData []models.StockDeficit

	err := repo.DB.
		Preload("Item.SupplierItems.Supplier").
		Where("parent_main_sku IN (?)", filter).
		Where("store_id = ?", store).
		Limit(pageSize).
		Offset(page).
		Find(&modelsData).Error

	return modelsData, err

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
		Preload("Item.SupplierItems", "supplier_id = ?", supplier).
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

type StockDeficitResult struct {
	ItemID  uint64  `gorm:"column:item_id"`
	Deficit float64 `gorm:"column:to_order"`
}

func (repo *StockDeficitImpl) StockDeficitByFatherOrder(father_id uint64) ([]StockDeficitResult, error) {
	var deficit []StockDeficitResult
	err1 := repo.DB.
		Table("order_lines AS ol").
		Select("ol.item_id , (ol.quantity  - IFNULL(ol2.quantity, 0) ) to_order ").
		Joins("LEFT join order_lines ol2 on ol2.item_id = ol.item_id and ol2.order_id = ol.order_id and ol2.store_id = 1").
		Where("ol.order_id in (select id from orders where father_order_id = ? )", father_id).
		Where("ol.store_id = 2").
		Where("(ol.quantity  - IFNULL(ol2.quantity, 0)) >0").
		Find(&deficit).Error

	if err1 != nil {
		fmt.Errorf("error al buscar registro existente: %w", err1)
		return nil, err1
	}
	return deficit, nil

}

func (repo *StockDeficitImpl) CallStockDefProc() {

	if err := repo.DB.Exec("CALL UpdateStockDeficitByStore();").Error; err != nil {
		log.Printf("Error ejecutando UpdateStockDeficitByStore: %v", err)
	}

}

func (repo *StockDeficitImpl) CallPendingStockProc() {
	// Llamada al segundo procedimiento almacenado
	if err := repo.DB.Exec("CALL UpdatePendingStocks();").Error; err != nil {
		log.Printf("Error ejecutando UpdatePendingStocks: %v", err)
	}

}
func (repo *StockDeficitImpl) CleanStockDeficit() error {
	err := repo.DB.
		Table("stock_deficits").
		Where("1 = 1").
		Updates(map[string]interface{}{
			"quantity":      0,
			"pending_stock": 0,
		}).Error

	if err != nil {
		log.Printf("Error limpiando StockDeficit: %v", err)
	}
	return err
}
