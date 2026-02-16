package domain

import "time"

// CreditStatus represents the state of a credit.
type CreditStatus string

const (
	CreditStatusEnCours  CreditStatus = "en_cours"
	CreditStatusEnRetard CreditStatus = "en_retard"
	CreditStatusPaye     CreditStatus = "paye"
)

// PaymentMethod represents how a payment was made.
type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "cash"
	PaymentCarte    PaymentMethod = "carte"
	PaymentVirement PaymentMethod = "virement"
)

// Payment represents a single payment against a credit.
type Payment struct {
	ID       string        `json:"id"`
	CreditID string        `json:"creditId"`
	Amount   float64       `json:"amount"`
	Date     time.Time     `json:"date"`
	Method   PaymentMethod `json:"method"`
}

// Credit represents money owed by a client for a sale.
type Credit struct {
	ID              string       `json:"id"`
	ClientID        string       `json:"clientId"`
	ClientName      string       `json:"clientName"`
	SaleID          string       `json:"saleId"`
	Amount          float64      `json:"amount"`
	RemainingAmount float64      `json:"remainingAmount"`
	Status          CreditStatus `json:"status"`
	CreatedAt       time.Time    `json:"createdAt"`
	DueDate         *time.Time   `json:"dueDate,omitempty"`
	Payments        []Payment    `json:"payments"`
}

// CreatePaymentRequest represents the payload to register a payment on a credit.
type CreatePaymentRequest struct {
	Amount float64       `json:"amount" validate:"required,gt=0"`
	Method PaymentMethod `json:"method" validate:"required"`
}
