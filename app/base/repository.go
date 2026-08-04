package base

import "gorm.io/gorm"

type Repository[T any, V any] interface {
	Name() string
	GetTableById(conn *gorm.DB, id string, preloads ...string) (T, error)
	GetUnscopedTableById(conn *gorm.DB, id string, preloads ...string) (T, error)
	GetViewById(conn *gorm.DB, id string, preloads ...string) (V, error)
	GetUnscopedViewById(conn *gorm.DB, id string, preloads ...string) (V, error)
	Create(conn *gorm.DB, data T) error
	Creates(conn *gorm.DB, data []T) error
	Update(conn *gorm.DB, data T) error
	Save(conn *gorm.DB, data T) error
	Delete(conn *gorm.DB, data T) error
}

type repositoryImpl[T any, V any] struct {
	name string
}

func NewRepository[T any, V any](name string) Repository[T, V] {
	return &repositoryImpl[T, V]{name: name}
}

func (r *repositoryImpl[T, V]) Name() string {
	return r.name
}

func (r *repositoryImpl[T, V]) applyPreloads(db *gorm.DB, preloads ...string) *gorm.DB {
	for _, p := range preloads {
		db = db.Preload(p)
	}
	return db
}

func (r *repositoryImpl[T, V]) GetTableById(conn *gorm.DB, id string, preloads ...string) (T, error) {
	var data T
	db := r.applyPreloads(conn, preloads...)
	err := db.First(&data, "id = ?", id).Error
	return data, err
}

func (r *repositoryImpl[T, V]) GetUnscopedTableById(conn *gorm.DB, id string, preloads ...string) (T, error) {
	var data T
	db := conn.Unscoped()
	db = r.applyPreloads(db, preloads...)
	err := db.First(&data, "id = ?", id).Error
	return data, err
}

func (r *repositoryImpl[T, V]) GetViewById(conn *gorm.DB, id string, preloads ...string) (V, error) {
	var data V
	db := r.applyPreloads(conn, preloads...)
	err := db.First(&data, "id = ?", id).Error
	return data, err
}

func (r *repositoryImpl[T, V]) GetUnscopedViewById(conn *gorm.DB, id string, preloads ...string) (V, error) {
	var data V
	db := conn.Unscoped()
	db = r.applyPreloads(db, preloads...)
	err := db.First(&data, "id = ?", id).Error
	return data, err
}

func (r *repositoryImpl[T, V]) Create(conn *gorm.DB, data T) error {
	return conn.Create(&data).Error
}

func (r *repositoryImpl[T, V]) Creates(conn *gorm.DB, data []T) error {
	return conn.Create(&data).Error
}

func (r *repositoryImpl[T, V]) Update(conn *gorm.DB, data T) error {
	return conn.Model(&data).Updates(data).Error
}

func (r *repositoryImpl[T, V]) Save(conn *gorm.DB, data T) error {
	return conn.Save(&data).Error
}

func (r *repositoryImpl[T, V]) Delete(conn *gorm.DB, data T) error {
	return conn.Delete(&data).Error
}
