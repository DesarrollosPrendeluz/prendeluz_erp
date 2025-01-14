package itemlocationrepo

import (
	"fmt"
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"

	"gorm.io/gorm"
)

type ItemLocationImpl struct {
	*repositories.GORMRepository[models.ItemLocation]
}

func NewInItemLocationRepository(db *gorm.DB) *ItemLocationImpl {
	return &ItemLocationImpl{repositories.NewGORMRepository(db, models.ItemLocation{})}
}

func (repo *ItemLocationImpl) FindByItemsAndLocation(mainSku string, location uint64) (models.ItemLocation, error) {
	var item models.ItemLocation
	result := repo.DB.Where("item_main_sku = ? and store_location_id = ?", mainSku, location).First(&item)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		fmt.Println("No se encontró el item y po eso se creará uno nuevo")
		item = models.ItemLocation{
			ItemMainSku:     mainSku,
			StoreLocationID: location,
			Stock:           0,
		}
		result = repo.DB.Create(&item)
	}
	return item, result.Error
}

// Busca un producto hijo en base a su aparición en la tabla parent_items
func (repo *ItemLocationImpl) FindByItemsAndStore(mainSku string, store uint64, pageSize int, offset int) ([]models.ItemLocation, error) {
	var item []models.ItemLocation
	subQuery := repo.DB.
		Table("store_locations").
		Select("id").
		Where("store_id = ?", store)
	result := repo.DB.Preload("StoreLocations").
		Where("item_main_sku = ? and store_location_id in (?)", mainSku, subQuery).
		Find(&item).
		Offset(offset).
		Limit(pageSize)
	return item, result.Error
}

// Busca un producto hijo en base a su aparición en la tabla parent_items
func (repo *ItemLocationImpl) FindByItem(mainSku string, pageSize int, offset int) ([]models.ItemLocation, error) {
	var item []models.ItemLocation

	result := repo.DB.Preload("StoreLocations.Store").
		Where("item_main_sku = ? ", mainSku).
		Find(&item).
		Order("store_location_id DESC").
		Offset(offset).
		Limit(pageSize)
	return item, result.Error
}

// Busca un producto hijo en base a su aparición en la tabla parent_items
func (repo *ItemLocationImpl) FindByItemAndLocation(mainSku string, locationId uint64) (models.ItemLocation, error) {
	var item models.ItemLocation

	err := repo.DB.
		Where("item_main_sku = ? ", mainSku).
		Where("store_location_id = ?", locationId).
		First(&item).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			modelCreate := models.ItemLocation{
				ItemMainSku:     mainSku,
				StoreLocationID: locationId,
				Stock:           0,
			}
			repo.DB.Create(&modelCreate)

		}
	}
	return item, err
}
