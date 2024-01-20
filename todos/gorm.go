package todos

import "gorm.io/gorm"

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (s *GormStore) NewTodo(todo *Todo) error {
	return s.db.Create(todo).Error
}

func (s *GormStore) GetTodos(todo *[]Todo) error {
	return s.db.Find(todo).Error
}

func (s *GormStore) DeleteTodo(todo *Todo, id int) error {
	return s.db.Delete(todo, id).Error
}
