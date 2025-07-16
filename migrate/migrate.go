package migrate

import (
	"inventory/initializers"
	"inventory/models"
	"log"
)

func init() {
	initializers.ConnectToDB()

}

func AutoMigrate() {
	log.Println("Running AutoMigrate...")

	initializers.DB.AutoMigrate(&models.Product{})
}
