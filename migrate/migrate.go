package migrate

import (
	"inventory/initializers"
	"inventory/models"
	"log"
)

func AutoMigrate() {
	log.Println("Running AutoMigrate...")
	err := initializers.DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
}
