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

// Retorna el producto de coincidir con el id de item así mismo retorna el proveedor y sus ubicaciones de alamcen
func (repo *ItemRepoImpl) FindByIdExtraData(id uint64) (models.Item, error) {
	var item models.Item

	result := repo.DB.Preload("SupplierItems.Supplier").Where("id LIKE ?", id).First(&item)

	return item, result.Error

}

func (repo *ItemRepoImpl) FindSonId(id uint64) (uint64, error) {
	var item models.Item
	var idChild uint64

	result := repo.DB.Preload("ChildRel").Where("id = ?", id).First(&item)
	fmt.Println("items type")
	fmt.Println(item.ItemType)
	if item.ItemType == "father" {
		fmt.Println("items clid")
		fmt.Printf("Items: %#v\n", item.ChildRel)
		idChild = item.ChildRel.ChildItemID
	} else {
		idChild = item.ID
	}
	fmt.Println("items")
	fmt.Printf("Items: %#v\n", item)
	return idChild, result.Error

}
