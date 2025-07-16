package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name" binding:"required"`     // Product name
	Description *string `json:"description"`                 // Optional description
	Price       float64 `json:"price" binding:"required"`    // Price
	Quantity    int     `json:"quantity" binding:"required"` // Quantity in stock
	SKU         *string `json:"sku"`                         // Optional Stock Keeping Unit
}
