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
	CountAll() (int64, error)
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

// Crea un nuevo regsitro del modelo referenciado
func (r *GORMRepository[T]) Create(item *T) error {
	result := r.DB.Create(item)
	return result.Error
}

//Crea un nuevos regsitros del modelo referenciado en el array

func (r *GORMRepository[T]) CreateAll(items *[]T) (int64, error) {
	result := r.DB.Create(items)
	return result.RowsAffected, result.Error
}

// Busca un regsitro de un modelo en base de datos por su id
func (r *GORMRepository[T]) FindByID(id uint64) (*T, error) {
	var item T
	result := r.DB.First(&item, id)
	return &item, result.Error
}

// Busca todas las ocurrencias de un modelo en base de datos pagiandas
func (r *GORMRepository[T]) FindAll(pageSize int, offset int) ([]T, error) {
	var items []T
	result := r.DB.Limit(pageSize).Offset(offset).Find(&items)
	return items, result.Error
}
func (r *GORMRepository[T]) FindAllWithOrder(pageSize int, offset int, orderBy string) ([]T, error) {
	var items []T
	result := r.DB.Limit(pageSize).Offset(offset).Order(orderBy).Find(&items)
	return items, result.Error
}

// Actualiza un registro en base de datos del modelo correspondiente
func (r *GORMRepository[T]) Update(item *T) error {
	return r.DB.Save(item).Error
}

// Actualiza todos los reistros de array
func (r *GORMRepository[T]) UpdateAll(items *[]T) error {
	return r.DB.Save(items).Error
}

// Borra la instancia de base de datos del modelo dado
func (r *GORMRepository[T]) Delete(id uint64) error {
	return r.DB.Delete(&r.model, id).Error
}

// Recuento de las instancias en base de datos
func (r *GORMRepository[T]) CountAll() (int64, error) {
	var count int64
	err := r.DB.Model(&r.model).Count(&count).Error
	return count, err
}
