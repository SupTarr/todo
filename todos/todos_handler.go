package todos

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/SupTarr/todo/my_context"
)

type TodoHandler struct {
	repo Repositorier
}

func NewTodoHandler(repo Repositorier) *TodoHandler {
	return &TodoHandler{repo: repo}
}

func (t *TodoHandler) GetTasks(c my_context.Context) {
	var todos []Todo

	err := t.repo.GetTodos(&todos)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) NewTask(c my_context.Context) {
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}

	if todo.Title == "sleep" {
		transactionID := c.TransactionID()
		aud := c.Audience()
		text := fmt.Sprintf("transction %s:, audience: %v not allowed", transactionID, aud)
		log.Println(text)
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": text,
		})
		return
	}

	err := t.repo.NewTodo(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"ID": todo.ID,
	})
}

func (t *TodoHandler) RemoveTask(c my_context.Context) {
	idParam := c.TodoID()
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}

	err = t.repo.DeleteTodo(&Todo{}, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error": err,
		})
		return
	}

	c.Status(http.StatusOK)
}
