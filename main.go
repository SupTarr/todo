package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/SupTarr/todo/auth"
	"github.com/SupTarr/todo/todos"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf(">> Please consider environment variables: %s\n", err)
	}

	db, err := gorm.Open(sqlite.Open("todos.sqlite"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&todos.Todo{})

	r := gin.Default()
	r.GET("/ping", PingPongHandler)

	r.GET("/token", auth.AccessToken("==signature=="))

	protected := r.Group("", auth.Protect([]byte("==signature==")))
	todoHandler := todos.NewTodoHandler(db)
	protected.POST("/todos", todoHandler.NewTask)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}

func PingPongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
