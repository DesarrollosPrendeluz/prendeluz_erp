package repositories

import (
	"prendeluz/erp/internal/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID         uint64 `gorm:"primaryKey"`
	Test_field string
	Created_At time.Time
}

func TestCreate(t *testing.T) {
	db, cleanup := SetUpTestDB(&TestModel{})

	defer cleanup()

	repo := repositories.NewGORMRepository(db, TestModel{})
	testModel := TestModel{Test_field: "test1"}
	err := repo.Create(&testModel)
	assert.NoError(t, err)
	assert.NotZero(t, testModel.ID)
}

func TestFind(t *testing.T) {
	db, cleanup := SetUpTestDB(&TestModel{})
	defer cleanup()

	repo := repositories.NewGORMRepository(db, TestModel{})
	testModel := TestModel{Test_field: "test1"}

	_, err := repo.FindByID(1)
	assert.Error(t, err)
	repo.Create(&testModel)

	testResult, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, testResult.ID, testModel.ID)
	assert.Equal(t, testResult.Test_field, testModel.Test_field)
	assert.Equal(t, testResult.Created_At, testModel.Created_At)

}

func TestFindAll(t *testing.T) {
	db, cleanup := SetUpTestDB(&TestModel{})
	defer cleanup()

	repo := repositories.NewGORMRepository(db, TestModel{})
	testModels := []TestModel{
		{Test_field: "test1"},
		{Test_field: "test2"},
		{Test_field: "test3"},
	}

	repo.CreateAll(&testModels)

	testResults, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, testResults, 3)
	assert.Contains(t, testResults, testModels[0])
	assert.Contains(t, testResults, testModels[1])
	assert.Contains(t, testResults, testModels[2])
}

func TestUpdate(t *testing.T) {

	db, cleanup := SetUpTestDB(&TestModel{})
	defer cleanup()

	repo := repositories.NewGORMRepository(db, TestModel{})
	testModels := []TestModel{
		{Test_field: "test1"},
		{Test_field: "test2"},
		{Test_field: "test3"},
	}

	repo.CreateAll(&testModels)

	updateModel, _ := repo.FindByID(2)
	updateModel.Test_field = "test_update"
	err := repo.Update(updateModel)

	expected, _ := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, expected, 3)
	assert.Contains(t, expected, testModels[0])
	assert.Contains(t, expected, *updateModel)
	assert.Contains(t, expected, testModels[2])
}

func TestDelete(t *testing.T) {

	db, cleanup := SetUpTestDB(&TestModel{})
	defer cleanup()

	repo := repositories.NewGORMRepository(db, TestModel{})
	testModels := []TestModel{
		{Test_field: "test1"},
		{Test_field: "test2"},
		{Test_field: "test3"},
	}

	repo.CreateAll(&testModels)

	err := repo.Delete(2)

	expected, _ := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, expected, 2)
	assert.Contains(t, expected, testModels[0])
	assert.Contains(t, expected, testModels[2])
	assert.NotContains(t, expected, testModels[1])
}
