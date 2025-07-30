package repository

import (
	"inventory/metrics"
	"inventory/models"

	"gorm.io/gorm"
)

type GormProductRepo struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) *GormProductRepo {
	return &GormProductRepo{DB: db}
}

func (gp *GormProductRepo) Create(product *models.Product) error {
	return metrics.ObserveDBQuery("create_product", func() error {
		return gp.DB.Create(product).Error
	})
}

func (gp *GormProductRepo) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := metrics.ObserveDBQuery("find_all_products", func() error {
		return gp.DB.Find(&products).Error
	})
	return products, err
}

func (gp *GormProductRepo) FindByID(id string) (*models.Product, error) {
	var product models.Product
	err := metrics.ObserveDBQuery("find_product_by_id", func() error {
		return gp.DB.First(&product, id).Error
	})
	return &product, err
}

func (gp *GormProductRepo) Update(product *models.Product) error {
	return metrics.ObserveDBQuery("update_product", func() error {
		return gp.DB.Save(product).Error
	})
}

func (gp *GormProductRepo) Delete(id string) error {
	err := metrics.ObserveDBQuery("delete_product", func() error {
		return gp.DB.Delete(&models.Product{}, id).Error
	})

	if err == nil {
		metrics.IncrementDeletedItems("products", 1)
	}

	return err
}
