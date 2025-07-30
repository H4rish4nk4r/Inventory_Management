package interfaces

import (
	"inventory/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll() ([]models.Product, error)
	FindByID(id string) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id string) error
}
