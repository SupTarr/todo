package repositories

import (
	"context"

	"github.com/SupTarr/todo/todos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBStore struct {
	*mongo.Collection
}

func NewMongoDBStore(db *mongo.Collection) *MongoDBStore {
	return &MongoDBStore{Collection: db}
}

func (s *MongoDBStore) NewTodo(todo *todos.Todo) error {
	_, err := s.Collection.InsertOne(context.Background(), todo)
	return err
}

func (s *MongoDBStore) GetTodos(todo *[]todos.Todo) error {
	_, err := s.Collection.Find(context.Background(), todo)
	return err
}

func (s *MongoDBStore) DeleteTodo(todo *todos.Todo, id int) error {
	filter := bson.D{{Key: "ID", Value: id}}
	_, err := s.Collection.DeleteOne(context.Background(), filter)
	return err
}
