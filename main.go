package main

import (
	"log"
	"os"

	docs "inventory/docs" 

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"inventory/controllers"
	"inventory/initializers"
	"inventory/migrate"
)

func init() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ Could not load .env file")
		} else {
			log.Println("✅ .env file loaded")
		}
	}

	log.Println("DB_HOST:", os.Getenv("DB_HOST"))
	log.Println("DB_USER:", os.Getenv("DB_USER"))
	log.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	log.Println("DB_NAME:", os.Getenv("DB_NAME"))
	log.Println("DB_PORT:", os.Getenv("DB_PORT"))

	initializers.ConnectToDB()
	migrate.AutoMigrate()
}

func main() {
	r := gin.Default()

	
	docs.SwaggerInfo.Title = "Inventory API"
	docs.SwaggerInfo.Description = "This is an API for managing inventory"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	
	r.POST("/products", controllers.CreateProduct)
	r.GET("/products", controllers.GetProducts)
	r.GET("/products/:id", controllers.GetProductByID)
	r.PUT("/products/:id", controllers.UpdateProduct)
	r.DELETE("/products/:id", controllers.DeleteProduct)

	if os.Getenv("APP_ENV") != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	r.Run() 
}
