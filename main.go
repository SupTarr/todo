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

type UserHandler struct {
	db *gorm.DB
}

func main() {
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&User{})

	db.Create(&User{Name: "Supakrit"})

	r := gin.Default()
	r.GET("/ping", PingPongHandler)

	userHandler := UserHandler{db: db}
	r.GET("/users", userHandler.GetUser)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func PingPongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var u User
	h.db.First(&u)
	c.JSON(200, u)
}
