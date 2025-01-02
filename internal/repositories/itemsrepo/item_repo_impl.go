package itemsrepo

import (
	"fmt"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type ItemRepoImpl struct {
	*repositories.GORMRepository[models.Item]
}

func NewItemRepository(db *gorm.DB) *ItemRepoImpl {
	return &ItemRepoImpl{repositories.NewGORMRepository(db, models.Item{})}
}

// Retorna el producto de coincidir el texto parcial o totalmente con el main_sku de la tabla items
func (repo *ItemRepoImpl) FindByMainSku(sku string) (models.Item, error) {
	var item models.Item

	result := repo.DB.Where("main_sku LIKE ?", "%"+sku+"%").First(&item)

	return item, result.Error

}

func (repo *ItemRepoImpl) FindByIdWithFatherPreload(id uint64) (models.Item, error) {
	var item models.Item

	result := repo.DB.Preload("FatherRel.Parent").Where("id = ?", id).First(&item)

	return item, result.Error

}

func (repo *ItemRepoImpl) FindByEan(ean string) ([]models.Item, error) {
	var item []models.Item

	result := repo.DB.Where("ean = ?", ean).Find(&item)

	return item, result.Error

}

// Retorna el producto de coincidir con el id de item así mismo retorna el proveedor y sus ubicaciones de alamcen
func (repo *ItemRepoImpl) FindByIdExtraData(id uint64) (models.Item, error) {
	var item models.Item

	result := repo.DB.Preload("SupplierItems.Supplier").Preload("ItemLocations.StoreLocations").Where("id LIKE ?", id).First(&item)

	return item, result.Error

}

func (repo *ItemRepoImpl) FindSonId(id uint64) (uint64, error) {
	var item models.Item
	var idChild uint64

	result := repo.DB.Preload("ChildRel").Where("id = ?", id).First(&item)

	if item.ItemType == "father" {

		idChild = item.ChildRel.Child.ID
	} else {
		idChild = item.ID
	}
	return idChild, result.Error

}

func (repo *ItemRepoImpl) FindByMainSkus(skus []string) (map[string]models.Item, error) {
	var items []models.Item
	skuMap := make(map[string]models.Item)

	result := repo.DB.Where("main_sku in ?", skus).Find(&items)

	for _, item := range items {
		skuMap[item.MainSKU] = item

	}

	return skuMap, result.Error

}

// Retorna el producto de coincidir el texto parcial o totalmente con el main_sku de la tabla items
func (repo *ItemRepoImpl) FindByFathersMainSkuOrEan(filter string) ([]models.Item, error) {
	var item []models.Item

	result := repo.DB.
		Select("main_sku").
		Where("item_type like 'father'").
		Where(repo.DB.Where("main_sku = ?", filter).Or("ean = ?", filter)).
		Find(&item)

	return item, result.Error
}

func (repo *ItemRepoImpl) FindByEanAndSupplierSku(ean string, supplier_sku string) ([]uint64, error) {
	var itemIds []uint64
	var results []models.Item

	if ean != "" || supplier_sku != "" {
		query := repo.DB.
			Table("items AS i").
			Select(" distinct i.id").
			Joins("inner join supplier_items as si on si.item_id = i.id")

		if ean != "" && supplier_sku != "" {
			// Ambos, ean y supplier_sku, tienen valor
			query = query.Where("i.ean = ? OR si.supplier_sku = ?", ean, supplier_sku)
		} else if ean != "" {
			// Solo ean tiene valor
			query = query.Where("i.ean = ?", ean)
		} else if supplier_sku != "" {
			// Solo supplier_sku tiene valor
			query = query.Where("si.supplier_sku = ?", supplier_sku)
		}

		errr := query.Find(&results).Error
		if errr != nil {
			fmt.Println(errr.Error())

		}

		itemIds = func() []uint64 {
			var ids []uint64
			for _, item := range results {
				ids = append(ids, item.ID)

			}
			return ids

		}()

	}
	return itemIds, nil

}
