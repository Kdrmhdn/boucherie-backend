package domain

import "time"

// Client represents a butcher shop customer.
type Client struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" validate:"required,min=2"`
	Phone       string    `json:"phone" validate:"required"`
	Email       string    `json:"email,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
	TotalCredit float64   `json:"totalCredit"`
	CreatedAt   time.Time `json:"createdAt"`
}

// CreateClientRequest represents the payload to create a new client.
type CreateClientRequest struct {
	Name  string `json:"name" validate:"required,min=2"`
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email,omitempty"`
}

// UpdateClientRequest represents the payload to update an existing client.
type UpdateClientRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Phone *string `json:"phone,omitempty"`
	Email *string `json:"email,omitempty"`
}
