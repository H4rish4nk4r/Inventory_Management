package migrate

import (
	"inventory/models"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	log.Println("Running AutoMigrate...")
	err := db.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
}
