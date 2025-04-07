package repositories

import (
	"github.com/SupTarr/todo/todos"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (s *GormStore) NewTodo(todo *todos.Todo) error {
	return s.db.Create(todo).Error
}

func (s *GormStore) GetTodos(todo *[]todos.Todo) error {
	return s.db.Find(todo).Error
}

func (s *GormStore) DeleteTodo(todo *todos.Todo, id int) error {
	return s.db.Delete(todo, id).Error
}
