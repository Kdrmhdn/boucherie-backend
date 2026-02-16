package domain

import "time"

// SaleItem represents a single line item in a sale.
type SaleItem struct {
	ID          string  `json:"id"`
	SaleID      string  `json:"saleId"`
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    float64 `json:"quantity"` // in kg
	Subtotal    float64 `json:"subtotal"`
}

// Sale represents a completed sale transaction.
type Sale struct {
	ID           string     `json:"id"`
	ClientID     string     `json:"clientId"`
	ClientName   string     `json:"clientName"`
	Items        []SaleItem `json:"items"`
	Total        float64    `json:"total"`
	PaidAmount   float64    `json:"paidAmount"`
	CreditAmount float64    `json:"creditAmount"`
	Date         time.Time  `json:"date"`
}

// CreateSaleItemRequest is used to add items when creating a sale.
type CreateSaleItemRequest struct {
	ProductID string  `json:"productId" validate:"required"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
}

// CreateSaleRequest represents the payload to register a new sale.
type CreateSaleRequest struct {
	ClientID   string                  `json:"clientId" validate:"required"`
	Items      []CreateSaleItemRequest `json:"items" validate:"required,min=1,dive"`
	PaidAmount float64                 `json:"paidAmount" validate:"gte=0"`
}
