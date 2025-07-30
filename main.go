package main

import (
	"log"
	"os"
	"time"

	"inventory/controllers"
	docs "inventory/docs"
	"inventory/initializers"
	"inventory/metrics" // <- you’ll need to create this package
	"inventory/middlewares"
	"inventory/migrate"
	"inventory/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func init() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Could not load .env file")
		} else {
			log.Println(".env file loaded")
		}
	}
}

func main() {
	r := gin.Default()

	// Setup DB connection
	var DB *gorm.DB
	DB = initializers.ConnectToDB()

	// Auto migrate DB schema
	migrate.AutoMigrate(DB)

	// Setup repository and controllers
	productRepo := repository.NewProductController(DB)
	productController := controllers.NewProductController(productRepo)

	// Swagger metadata
	docs.SwaggerInfo.Title = "Inventory API"
	docs.SwaggerInfo.Description = "This is an API for managing inventory"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Register middleware
	r.Use(middlewares.PrometheusMiddleware())

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	r.POST("/api/v1/products", productController.CreateProduct)
	r.GET("/api/v1/products", productController.GetProducts)
	r.GET("/api/v1/products/:id", productController.GetProductByID)
	r.PUT("/api/v1/products", productController.UpdateProduct)
	r.DELETE("/api/v1/products/:id", productController.DeleteProduct)

	// Swagger docs
	if os.Getenv("APP_ENV") != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// Start background goroutine to collect DB metrics
	go func() {
		for {
			sqlDB, err := DB.DB()
			if err != nil {
				log.Fatalf("Failed to get *sql.DB: %v", err)
			}
			metrics.UpdateTableRowCount(sqlDB, "products")

			time.Sleep(30 * time.Second) // adjust interval as needed
		}
	}()

	// Run server
	r.Run()
}
