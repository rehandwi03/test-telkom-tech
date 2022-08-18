package response

import (
	"github.com/google/uuid"
)

type CartResponse struct {
	ID       uuid.UUID             `json:"id"`
	Products []CartResponseProduct `json:"products"`
	FullName string                `json:"full_name"`
}

type CartResponseProduct struct {
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
}

type CartProduct struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Price             float64   `json:"price"`
	IsDiscount        bool      `json:"is_discount"`
	StartDateDiscount string    `json:"start_date_discount"`
	EndDateDiscount   string    `json:"end_date_discount"`
}
