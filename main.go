package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

func main() {
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&User{})

	db.Create(&User{Name: "Supakrit"})

	r := gin.Default()
	r.GET("/ping", pingPongHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func pingPongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
