package request

import (
	"github.com/google/uuid"
)

type CartAddRequest struct {
	FullName string                `json:"full_name" binding:"required"`
	Product  CartAddProductRequest `json:"product" binding:"required"`
}

type CartAddProductRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
}

type CartCriteria struct {
	FullName    string `json:"full_name"`
	Quantity    string `json:"quantity"`
	ProductName string `json:"product_name"`
}
