package port

import (
	"boucherie-api/internal/domain"
	"context"
)

// ClientRepository defines the contract for client persistence.
type ClientRepository interface {
	FindAll(ctx context.Context) ([]domain.Client, error)
	FindByID(ctx context.Context, id string) (*domain.Client, error)
	Create(ctx context.Context, client *domain.Client) error
	Update(ctx context.Context, client *domain.Client) error
	UpdateTotalCredit(ctx context.Context, clientID string, delta float64) error
}

// ProductRepository defines the contract for product persistence.
type ProductRepository interface {
	FindAll(ctx context.Context, category *domain.MeatCategory) ([]domain.Product, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error
}

// SaleRepository defines the contract for sale persistence.
type SaleRepository interface {
	FindAll(ctx context.Context, clientID *string, date *string) ([]domain.Sale, error)
	FindByID(ctx context.Context, id string) (*domain.Sale, error)
	Create(ctx context.Context, sale *domain.Sale) error
}

// CreditRepository defines the contract for credit persistence.
type CreditRepository interface {
	FindAll(ctx context.Context, status *domain.CreditStatus) ([]domain.Credit, error)
	FindByID(ctx context.Context, id string) (*domain.Credit, error)
	FindByClientID(ctx context.Context, clientID string) ([]domain.Credit, error)
	Create(ctx context.Context, credit *domain.Credit) error
	Update(ctx context.Context, credit *domain.Credit) error
	AddPayment(ctx context.Context, payment *domain.Payment) error
}

// OrderRepository defines the contract for order persistence.
type OrderRepository interface {
	FindAll(ctx context.Context, status *domain.OrderStatus) ([]domain.Order, error)
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	Create(ctx context.Context, order *domain.Order) error
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id string) error
}

// DashboardStats holds aggregated data for the dashboard.
type DashboardStats struct {
	TotalRevenue    float64        `json:"totalRevenue"`
	TotalCash       float64        `json:"totalCash"`
	TotalCredit     float64        `json:"totalCredit"`
	AvgTicket       float64        `json:"avgTicket"`
	SalesCount      int            `json:"salesCount"`
	ClientsServed   int            `json:"clientsServed"`
	TopDebtors      []domain.Client `json:"topDebtors"`
}
