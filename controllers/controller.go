package controllers

import (
	"inventory/initializers"
	"inventory/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create a new product with input payload
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Product to create"
// @Success      201      {object}  models.Product
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /products [post]
func CreateProduct(c *gin.Context) {
	var input models.Product

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
		SKU:         input.SKU,
	}

	if result := initializers.DB.Create(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProducts godoc
// @Summary      Get all products
// @Description  Retrieve all products from the database
// @Tags         products
// @Produce      json
// @Success      200  {array}   models.Product
// @Failure      500  {object}  map[string]string
// @Router       /products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product

	if result := initializers.DB.Find(&products); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary      Get a product by ID
// @Description  Retrieve a product by its ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.Product
// @Failure      404  {object}  map[string]string
// @Router       /products/{id} [get]
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if result := initializers.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update product details by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int            true  "Product ID"
// @Param        product  body      models.Product  true  "Updated product data"
// @Success      200      {object}  models.Product
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if result := initializers.DB.First(&product, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Quantity = input.Quantity
	product.SKU = input.SKU

	if result := initializers.DB.Save(&product); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product by its ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if result := initializers.DB.Delete(&models.Product{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
