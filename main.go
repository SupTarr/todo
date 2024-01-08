package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/SupTarr/todo/todos"
)

func main() {
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&todos.Todo{})

	r := gin.Default()
	r.GET("/ping", PingPongHandler)

	todoHandler := todos.NewTodoHandler(db)
	r.POST("/todos", todoHandler.NewTask)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func PingPongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
