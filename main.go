package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"product-service/api"
	"product-service/product"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func run() error {

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Database connection
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	// Dependency initialization
	repository := product.NewRepository(db)
	service := product.NewService(repository)
	handler := product.NewHandler(service)

	// HTTP server
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	api.RegisterHandlers(router, handler)

	return router.Run(":8080")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
