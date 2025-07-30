package controllers

import (
	"inventory/interfaces"
	"inventory/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// type DBI interface {
// 	Create(models.Product) error
// }

// type ProdutController struct {
// 	DB DBI
// }

// func NewProductController(db DBI) ProdutController {
// 	return ProdutController{DB: db}
// }

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
// func (pc ProdutController) CreateProduct(c *gin.Context)

type ProductController struct {
	Repo interfaces.ProductRepository
}

func NewProductController(repo interfaces.ProductRepository) *ProductController {
	return &ProductController{Repo: repo}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var input models.Product
	err := c.ShouldBindJSON(&input)

	if err != nil {
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

	if err := pc.Repo.Create(&product); err != nil {
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
func (pc *ProductController) GetProducts(c *gin.Context) {
	products, err := pc.Repo.FindAll()

	if err != nil {
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
func (pc *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if _, err := pc.Repo.FindByID(id); err != nil {
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
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// Step 1: Find existing product
	product, err := pc.Repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Step 2: Parse new values
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 3: Update the fields
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Quantity = input.Quantity
	product.SKU = input.SKU

	// Step 4: Save to DB
	if err := pc.Repo.Update(product); err != nil {
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
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if result := pc.Repo.Delete(id); result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
