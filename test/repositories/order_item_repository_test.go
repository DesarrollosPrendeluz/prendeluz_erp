package repositories

import (
	"prendeluz/erp/internal/models"
	"prendeluz/erp/internal/repositories"
	"prendeluz/erp/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindByOrderCode(t *testing.T) {
	db, cleanup := test.SetUpTestDB(&models.OrderItem{})
	defer cleanup()

	repo := repositories.NewOrderItemRepository(db)
	testModels := []models.OrderItem{
		{ItemID: 1, OrderID: 1, Amount: 5},
		{ItemID: 2, OrderID: 1, Amount: 15},
		{ItemID: 1, OrderID: 2, Amount: 45},
		{ItemID: 5, OrderID: 5, Amount: 2},
	}

	repo.CreateAll(&testModels)

	testResults, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, testResults, 4)
	assert.Contains(t, testResults, testModels[0])
	assert.Contains(t, testResults, testModels[1])
	assert.Contains(t, testResults, testModels[2])
	assert.Contains(t, testResults, testModels[3])
}
