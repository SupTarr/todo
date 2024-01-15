package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/SupTarr/todo/auth"
	"github.com/SupTarr/todo/todos"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

var limiter = rate.NewLimiter(5, 5)

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Printf(">> Please consider environment variables: %s\n", err)
	}

	db, err := gorm.Open(sqlite.Open("todos.sqlite"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&todos.Todo{})

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})
	r.GET("/pingz", PingPongHandler)
	r.GET("/limitz", LimitedHandler)
	r.GET("/x", XHandler)

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))
	todoHandler := todos.NewTodoHandler(db)
	protected.POST("/todos", todoHandler.NewTask)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func PingPongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func XHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"buildcommit": buildcommit, "buildtime": buildtime})
}

func LimitedHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
