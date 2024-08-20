package repositories

import (
	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(item *T) error
	FindByID(id uint64) (*T, error)
	FindAll() ([]T, error)
	Update(item *T) error
	Delete(id uint64) error
}

type GORMRepository[T any] struct {
	db    *gorm.DB
	model T
}

func NewGORMRepository[T any](db *gorm.DB, model T) *GORMRepository[T] {
	return &GORMRepository[T]{db: db, model: model}
}

func (r *GORMRepository[T]) Create(item *T) error {
	result := r.db.Create(item)
	return result.Error
}

func (r *GORMRepository[T]) CreateAll(items *[]T) (int64, error) {
	result := r.db.Create(items)
	return result.RowsAffected, result.Error
}
func (r *GORMRepository[T]) FindByID(id uint64) (*T, error) {
	var item T
	result := r.db.First(&item, id)
	return &item, result.Error
}

func (r *GORMRepository[T]) FindAll() ([]T, error) {
	var items []T
	result := r.db.Find(&items)
	return items, result.Error
}

func (r *GORMRepository[T]) Update(item *T) error {
	return r.db.Save(item).Error
}

func (r *GORMRepository[T]) Delete(id uint64) error {
	return r.db.Delete(&r.model, id).Error
}
