package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/SupTarr/todo/auth"
	"github.com/SupTarr/todo/my_context"
	"github.com/SupTarr/todo/repositories"
	"github.com/SupTarr/todo/router"
	"github.com/SupTarr/todo/todos"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

var limiter = rate.NewLimiter(5, 5)

func initConfig() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("Please consider environment variables: %s\n", err)
	}

	if os.Getenv("SIGN") == "" {
		log.Fatal("SIGN environment variable is not set")
	}
}

func initGin() *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8080",
	}

	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"TransactionID",
	}

	r.Use(cors.New(config))
	return r
}

func initGorm(dns string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	db.AutoMigrate(&todos.Todo{})
	return db
}

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove("/tmp/live")

	initConfig()

	dns := os.Getenv("DB_CONN")
	db := initGorm(dns)

	r := initGin()
	gormRepo := repositories.NewGormStore(db)
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})

	r.GET("/pingz", router.NewGinHandler(PingPongHandler))
	r.GET("/limitz", router.NewGinHandler(LimitedHandler))
	r.GET("/x", router.NewGinHandler(XHandler))
	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))
	todoHandler := todos.NewTodoHandler(gormRepo)
	protected.POST("/todos", router.NewGinHandler(todoHandler.NewTask))
	protected.GET("/todos", router.NewGinHandler(todoHandler.GetTasks))
	protected.DELETE("/todos/:id", router.NewGinHandler(todoHandler.RemoveTask))

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

	gracefullyShutdown(ctx, srv)
}

func gracefullyShutdown(ctx context.Context, srv *http.Server) {
	<-ctx.Done()

	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func PingPongHandler(c my_context.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func XHandler(c my_context.Context) {
	c.JSON(http.StatusOK, gin.H{"buildcommit": buildcommit, "buildtime": buildtime})
}

func LimitedHandler(c my_context.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
