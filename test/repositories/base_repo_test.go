package repositories

import (
	"prendeluz/erp/internal/repositories"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestModel struct {
	ID         uint64 `gorm:"primaryKey"`
	Test_field string
	Created_At time.Time `gorm:"autoCreateTime"`
}

func TestCreate(t *testing.T) {

	mockDb, mock, _ := sqlmock.New()
	db, _ := gorm.Open(mysql.New(mysql.Config{Conn: mockDb, SkipInitializeWithVersion: true}), &gorm.Config{})
	repo := repositories.NewGORMRepository(db, TestModel{})
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `test_models` \\(`test_field`,`created_at`\\) VALUES \\(\\?\\,\\?\\)").WithArgs("test1", sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	testModel := TestModel{Test_field: "test1"}
	repo.Create(&testModel)

	err := mock.ExpectationsWereMet()

	assert.NoError(t, err)
	assert.NotZero(t, testModel.ID)
}

func TestFind(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	db, _ := gorm.Open(mysql.New(mysql.Config{Conn: mockDb, SkipInitializeWithVersion: true}), &gorm.Config{})
	repo := repositories.NewGORMRepository(db, TestModel{})

	rows := sqlmock.NewRows([]string{"id", "test_field", "created_at"}).AddRow(1, "test1", time.Now())

	mock.ExpectQuery("SELECT \\* FROM `test_models` WHERE `test_models`.`id` = \\? ORDER BY `test_models`.`id` LIMIT \\?").WithArgs(1, 1).WillReturnRows(rows)
	testResult, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, testResult.ID, uint64(1))
	assert.Equal(t, testResult.Test_field, "test1")

	err = mock.ExpectationsWereMet()

	assert.NoError(t, err)

}

func TestFindAll(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	db, _ := gorm.Open(mysql.New(mysql.Config{Conn: mockDb, SkipInitializeWithVersion: true}), &gorm.Config{})
	repo := repositories.NewGORMRepository(db, TestModel{})

	rows := sqlmock.NewRows([]string{"id", "test_field", "created_at"}).AddRow(1, "test1", time.Now()).AddRow(1, "test2", time.Now()).AddRow(1, "test2", time.Now())

	mock.ExpectQuery("SELECT \\* FROM `test_models`  LIMIT \\?").WithArgs(10).WillReturnRows(rows)
	testResults, err := repo.FindAll(10, 0)

	assert.NoError(t, err)
	mock.ExpectationsWereMet()
	assert.Len(t, testResults, 3)
}

func TestUpdate(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	db, _ := gorm.Open(mysql.New(mysql.Config{Conn: mockDb, SkipInitializeWithVersion: true}), &gorm.Config{})
	repo := repositories.NewGORMRepository(db, TestModel{})

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE `+"`test_models`"+` SET `+"`test_field`=\\?,`created_at`=\\? WHERE `id` = \\?").WithArgs("new_value", sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
	updateModel := TestModel{ID: 1, Test_field: "new_value", Created_At: time.Now()}
	mock.ExpectCommit()
	err := repo.Update(&updateModel)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// func TestDelete(t *testing.T) {
//
// 	repo := repositories.NewGORMRepository(db, TestModel{})
// 	testModels := []TestModel{
// 		{Test_field: "test1"},
// 		{Test_field: "test2"},
// 		{Test_field: "test3"},
// 	}
//
// 	repo.CreateAll(&testModels)
//
// 	err := repo.Delete(2)
//
// 	expected, _ := repo.FindAll(10, 0)
// 	assert.NoError(t, err)
// 	assert.Len(t, expected, 2)
// 	assert.Contains(t, expected, testModels[0])
// 	assert.Contains(t, expected, testModels[2])
// 	assert.NotContains(t, expected, testModels[1])
// }
