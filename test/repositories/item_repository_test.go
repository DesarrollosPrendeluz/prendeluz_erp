package repositories

// import (
// 	"prendeluz/erp/internal/models"
// 	"prendeluz/erp/internal/repositories/itemsrepo"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// )

// func TestFindByMainSku(t *testing.T) {
// 	db, cleanup := SetUpTestDB(&models.Item{}, &models.Category{}, &models.CategoryStatusType{}, &models.TypeOfCategories{})
//
// 	defer cleanup()
// 	name := "Test item"
//
// 	repo := itemsrepo.NewItemRepository(db)
// 	testItem := &models.Item{
// 		MainSKU:        "SKU123",
// 		Name:           &name,
// 		ItemType:       "father",
// 		AssignmentCost: 10.0,
// 	}
//
// 	err := repo.Create(testItem)
//
// 	assert.NoError(t, err)
// 	assert.NotZero(t, testItem.ID)
//
// }
