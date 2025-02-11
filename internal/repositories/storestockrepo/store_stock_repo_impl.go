package storestockrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"strconv"

	"gorm.io/gorm"
)

type StoreStockRepoImpl struct {
	*repositories.GORMRepository[models.StoreStock]
}

func NewStoreStockRepository(db *gorm.DB) *StoreStockRepoImpl {
	return &StoreStockRepoImpl{repositories.NewGORMRepository(db, models.StoreStock{})}
}

func (repo *StoreStockRepoImpl) FindByItem(sku_parent string) (models.StoreStock, error) {
	var storeStocks models.StoreStock

	results := repo.DB.Where("parent_sku LIKE ?", "%"+sku_parent+"%").First(&storeStocks)

	return storeStocks, results.Error
}

func (repo *StoreStockRepoImpl) FindByItemAndStore(sku_parent string, store_id string) (models.StoreStock, error) {
	var storeStocks models.StoreStock

	results := repo.DB.Where("parent_main_sku LIKE ? AND store_id = ? ", "%"+sku_parent+"%", store_id).First(&storeStocks)
	storeId, _ := strconv.ParseUint(store_id, 10, 64)
	if results.Error == gorm.ErrRecordNotFound {
		storeStock := models.StoreStock{
			SKU_Parent:     sku_parent,
			StoreID:        storeId,
			Amount:         0,
			ReservedAmount: 0,
		}
		if err := repo.DB.Create(&storeStock).Error; err != nil {
			return storeStock, err
		}
		return storeStock, nil
	}

	return storeStocks, results.Error
}

func (repo *StoreStockRepoImpl) FindByStore(idStore uint64, pageSize int, offset int) ([]models.StoreStock, error) {
	var storeStocks []models.StoreStock

	results := repo.DB.Limit(pageSize).Offset(offset).Preload("Item.AsinRel").Where("store_id = ?", idStore).Find(&storeStocks)

	return storeStocks, results.Error

}

type ItemsLocations struct {
	ItemSku string
	Ean     string
	Stock   int
	Code    string
	//ItemLocations []models.ItemLocation
}

type StoreAndlocations struct {
	StoreStocks   []models.StoreStock
	ItemsLocation []ItemsLocations
}

func (repo *StoreStockRepoImpl) FindByStoreWithLocations(idStore uint64) (StoreAndlocations, error) {
	var storeStocks []models.StoreStock
	var storeStock []models.StoreStock
	var stocks StoreAndlocations
	batchSize := 1000
	data := repo.DB.
		Preload("Item").
		Preload("Locations.StoreLocations").
		Where("store_stocks.store_id = ?", idStore).
		FindInBatches(&storeStock, batchSize, func(tx *gorm.DB, batch int) error {
			storeStocks = append(storeStocks, storeStock...)
			return nil // Continuar con el siguiente lote
		})
	stocks.StoreStocks = storeStocks
	for _, stock := range stocks.StoreStocks {
		var newLoc ItemsLocations
		newLoc.ItemSku = stock.SKU_Parent
		newLoc.Ean = stock.Item.EAN
		if len(*stock.Locations) > 0 {
			for _, loc := range *stock.Locations {
				if loc.StoreLocations.StoreID == idStore {
					newLoc.Code = loc.StoreLocations.Code
					newLoc.Stock = loc.Stock
					stocks.ItemsLocation = append(stocks.ItemsLocation, newLoc)

				}

			}

		}

	}

	return stocks, data.Error

}

func (repo *StoreStockRepoImpl) FindByStoreAndSearchParams(idStore uint64, searchParam string, pageSize int, offset int) ([]models.StoreStock, error) {
	var storeStocks []models.StoreStock
	var itemsParent []interface{}
	var results []map[string]interface{}

	// Ejecutar la consulta y almacenar los resultados en el slice de mapas
	repo.DB.
		Select("IF(items.item_type = 'father', items.main_sku, items_parent_ref.main_sku) AS father_skus").
		Limit(pageSize).
		Offset(offset).
		Joins("left JOIN asins ON asins.item_id = items.id").
		Joins("left JOIN supplier_items ON supplier_items.item_id = items.id").
		Joins("Left JOIN item_parents ON item_parents.child_item_id = items.id").
		Joins("Left JOIN items as items_parent_ref ON item_parents.parent_item_id = items_parent_ref.id").
		Where("items.main_sku LIKE ? OR items.ean LIKE ? OR asins.code LIKE ? OR asins.ean LIKE ? OR supplier_items.supplier_sku LIKE ?",
			"%"+searchParam+"%", "%"+searchParam+"%", "%"+searchParam+"%", "%"+searchParam+"%", "%"+searchParam+"%").
		Table("items").
		Find(&results)

	// Iterar sobre los resultados y acceder a los datos directamente
	for _, row := range results {
		if parentItemID, ok := row["father_skus"]; ok { // Comprobar si "parent_item_id" existe en el resultado
			itemsParent = append(itemsParent, parentItemID)
		}
	}

	resultsQuery := repo.DB.Limit(pageSize).Offset(offset).Preload("Item").Where("store_id = ?", idStore).Where("parent_main_sku in ?", itemsParent).Find(&storeStocks)
	//fmt.Printf(json.Encoder(storeStocks))

	return storeStocks, resultsQuery.Error

}
