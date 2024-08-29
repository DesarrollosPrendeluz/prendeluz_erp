package repositories

import (
	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(item *T) error
	CreateAll(items *[]T) (int64, error)
	FindByID(id uint64) (*T, error)
	FindAll(pagesize int, offset int) ([]T, error)
	Update(item *T) error
	UpdateAll(items *[]T)
	Delete(id uint64) error
	SetDB(db *gorm.DB)
}

type GORMRepository[T any] struct {
	DB    *gorm.DB
	model T
}

func NewGORMRepository[T any](db *gorm.DB, model T) *GORMRepository[T] {
	return &GORMRepository[T]{DB: db, model: model}
}
func (r *GORMRepository[T]) SetDB(db *gorm.DB) {
	r.DB = db
}
func (r *GORMRepository[T]) Create(item *T) error {
	result := r.DB.Create(item)
	return result.Error
}

func (r *GORMRepository[T]) CreateAll(items *[]T) (int64, error) {
	result := r.DB.Create(items)
	return result.RowsAffected, result.Error
}
func (r *GORMRepository[T]) FindByID(id uint64) (*T, error) {
	var item T
	result := r.DB.First(&item, id)
	return &item, result.Error
}

func (r *GORMRepository[T]) FindAll(pageSize int, offset int) ([]T, error) {
	var items []T
	result := r.DB.Limit(pageSize).Offset(offset).Find(&items)
	return items, result.Error
}

func (r *GORMRepository[T]) Update(item *T) error {
	return r.DB.Save(item).Error
}
func (r *GORMRepository[T]) UpdateAll(items *[]T) error {
	return r.DB.Save(items).Error
}

func (r *GORMRepository[T]) Delete(id uint64) error {
	return r.DB.Delete(&r.model, id).Error
}
