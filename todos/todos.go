package todos

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title string `json:"text" binding:"required"`
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

func (t *TodoHandler) GetTasks(c *gin.Context) {
	var todos []Todo

	err := t.store.GetTodos(&todos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) NewTask(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if todo.Title == "sleep" {
		transactionID := c.Request.Header.Get("TransactionID")
		aud, _ := c.Get("aud")
		text := fmt.Sprintf("transction %s:, audience: %v not allowed", transactionID, aud)
		log.Println(text)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": text,
		})
		return
	}

	err := t.store.NewTodo(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}

func (t *TodoHandler) RemoveTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = t.store.DeleteTodo(&Todo{}, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.Status(http.StatusOK)
}
