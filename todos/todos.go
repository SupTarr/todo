package todos

import (
	"log"
	"net/http"

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

type TodoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (t *TodoHandler) GetTasks(c *gin.Context) {
	var todo []Todo

	r := t.db.Find(&todo)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, todo)
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
		log.Println(transactionID, "not allowed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "not allowed",
		})
		return
	}

	r := t.db.Create(&todo)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}
