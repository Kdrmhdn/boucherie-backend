package domain

import "time"

// OrderStatus represents the lifecycle state of an order.
type OrderStatus string

const (
	OrderStatusEnAttente OrderStatus = "en_attente"
	OrderStatusConfirmee OrderStatus = "confirmee"
	OrderStatusPrete     OrderStatus = "prete"
	OrderStatusLivree    OrderStatus = "livree"
	OrderStatusAnnulee   OrderStatus = "annulee"
)

// OrderItem represents a single product line in an order.
type OrderItem struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"orderId"`
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    float64 `json:"quantity"` // in kg
}

// Order represents a customer pre-order or reservation.
type Order struct {
	ID          string      `json:"id"`
	ClientID    string      `json:"clientId"`
	ClientName  string      `json:"clientName"`
	ClientPhone string      `json:"clientPhone"`
	Items       []OrderItem `json:"items"`
	PickupDate  time.Time   `json:"pickupDate"`
	Notes       string      `json:"notes,omitempty"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"createdAt"`
}

// CreateOrderItemRequest is used when creating an order.
type CreateOrderItemRequest struct {
	ProductID string  `json:"productId" validate:"required"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
}

// CreateOrderRequest represents the payload to create a new order.
type CreateOrderRequest struct {
	ClientID   string                   `json:"clientId" validate:"required"`
	Items      []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
	PickupDate string                   `json:"pickupDate" validate:"required"`
	Notes      string                   `json:"notes,omitempty"`
}

// UpdateOrderRequest represents the payload to update an order.
type UpdateOrderRequest struct {
	Status     *OrderStatus `json:"status,omitempty"`
	PickupDate *string      `json:"pickupDate,omitempty"`
	Notes      *string      `json:"notes,omitempty"`
}
