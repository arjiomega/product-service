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
)

func run() error {

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return router.Run(":" + port)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
