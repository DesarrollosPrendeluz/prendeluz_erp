package storestockrepo

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

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

func (repo *StoreStockRepoImpl) FindByStore(idStore uint64, pageSize int, offset int) ([]models.StoreStock, error) {
	var storeStocks []models.StoreStock

	results := repo.DB.Limit(pageSize).Offset(offset).Preload("Item").Where("store_id = ?", idStore).Find(&storeStocks)

	return storeStocks, results.Error

}

func (repo *StoreStockRepoImpl) FindByStoreAndSearchParams(idStore uint64, searchParam string, pageSize int, offset int) ([]models.StoreStock, error) {
	var storeStocks []models.StoreStock
	//TODO:filtros avanzados de busqueda
	// var itemsParent []interface{}
	// var results []map[string]interface{}

	// // Ejecutar la consulta y almacenar los resultados en el slice de mapas
	// repo.DB.
	// 	Select("IF(items.item_type == father, items.main_sku, (SELECT some_field FROM some_table WHERE ...)) AS father_skus").
	// 	Limit(pageSize).
	// 	Offset(offset).
	// 	Joins("JOIN asins ON asins.item_id = items.id").
	// 	Joins("JOIN supplier_items ON supplier_items.item_id = items.id").
	// 	//Joins("JOIN item_parents ON item_parents.child_item_id = items.id").
	// 	Where("items.main_sku LIKE ? OR items.ean LIKE ? OR asins.code LIKE ? OR asins.ean LIKE ? OR supplier_items.supplier_sku LIKE ?",
	// 		"%"+searchParam+"%", "%"+searchParam+"%", "%"+searchParam+"%", "%"+searchParam+"%", "%"+searchParam+"%").
	// 	Table("items").
	// 	Find(&results)

	// // Iterar sobre los resultados y acceder a los datos directamente
	// for _, row := range results {
	// 	if parentItemID, ok := row["father_skus"]; ok { // Comprobar si "parent_item_id" existe en el resultado
	// 		itemsParent = append(itemsParent, parentItemID)
	// 	}
	// }
	// fmt.Printf("Resulting StoreStocks: %+v\n", itemsParent)

	resultsQuery := repo.DB.Limit(pageSize).Offset(offset).Preload("Item").Where("store_id = ?", idStore).Find(&storeStocks)
	//fmt.Printf(json.Encoder(storeStocks))

	return storeStocks, resultsQuery.Error

}
