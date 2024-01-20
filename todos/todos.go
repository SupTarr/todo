package todos

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Todo struct {
	ID        uint   `gorm:"primarykey"`
	Title     string `json:"text" binding:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Todo) TableName() string {
	return "todos"
}

type storer interface {
	GetTodos(*[]Todo) error
	NewTodo(*Todo) error
	DeleteTodo(*Todo, int) error
}

type TodoHandler struct {
	store storer
}

func NewTodoHandler(store storer) *TodoHandler {
	return &TodoHandler{store: store}
}

type Context interface {
	Bind(interface{}) error
	TodoID() string
	TransactionID() string
	Audience() string
	Status(int)
	JSON(int, interface{})
}

func (t *TodoHandler) GetTasks(c Context) {
	var todos []Todo

	err := t.store.GetTodos(&todos)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) NewTask(c Context) {
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if todo.Title == "sleep" {
		transactionID := c.TransactionID()
		aud := c.Audience()
		text := fmt.Sprintf("transction %s:, audience: %v not allowed", transactionID, aud)
		log.Println(text)
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": text,
		})
		return
	}

	err := t.store.NewTodo(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"ID": todo.ID,
	})
}

func (t *TodoHandler) RemoveTask(c Context) {
	idParam := c.TodoID()
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	err = t.store.DeleteTodo(&Todo{}, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err,
		})
		return
	}

	c.Status(http.StatusOK)
}
