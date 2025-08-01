package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"inventory/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockRepo is a mock implementation of the ProductRepository interface
type MockRepo struct {
	CreateFn   func(product *models.Product) error
	FindAllFn  func() ([]models.Product, error)
	FindByIDFn func(id string) (*models.Product, error)
	UpdateFn   func(product *models.Product) error
	DeleteFn   func(id string) error
}

func (m *MockRepo) Create(product *models.Product) error        { return m.CreateFn(product) }
func (m *MockRepo) FindAll() ([]models.Product, error)          { return m.FindAllFn() }
func (m *MockRepo) FindByID(id string) (*models.Product, error) { return m.FindByIDFn(id) }
func (m *MockRepo) Update(product *models.Product) error        { return m.UpdateFn(product) }
func (m *MockRepo) Delete(id string) error                      { return m.DeleteFn(id) }

// setupRouter initializes a test Gin router with the given controller
func setupRouter(controller *ProductController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/products", controller.CreateProduct)
	return r
}

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := &MockRepo{
		CreateFn: func(product *models.Product) error {
			return nil
		},
	}

	controller := NewProductController(mockRepo)
	router := setupRouter(controller)

	product := models.Product{
		Name:     "Mocked Product",
		Price:    49.99,
		Quantity: 10,
	}
	body, _ := json.Marshal(product)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateProduct_InvalidInput(t *testing.T) {
	mockRepo := &MockRepo{}

	controller := NewProductController(mockRepo)
	router := setupRouter(controller)

	invalidJSON := []byte(`{}`) // missing required fields

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateProduct_RepoFails(t *testing.T) {
	mockRepo := &MockRepo{
		CreateFn: func(product *models.Product) error {
			return errors.New("mock DB failure")
		},
	}

	controller := NewProductController(mockRepo)
	router := setupRouter(controller)

	product := models.Product{
		Name:     "Test Product",
		Price:    20.00,
		Quantity: 5,
	}
	body, _ := json.Marshal(product)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
