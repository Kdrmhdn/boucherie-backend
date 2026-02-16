package service

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/port"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// SaleService handles sale business logic.
type SaleService struct {
	saleRepo    port.SaleRepository
	productRepo port.ProductRepository
	clientRepo  port.ClientRepository
	creditRepo  port.CreditRepository
}

// NewSaleService creates a new sale service.
func NewSaleService(
	saleRepo port.SaleRepository,
	productRepo port.ProductRepository,
	clientRepo port.ClientRepository,
	creditRepo port.CreditRepository,
) *SaleService {
	return &SaleService{
		saleRepo:    saleRepo,
		productRepo: productRepo,
		clientRepo:  clientRepo,
		creditRepo:  creditRepo,
	}
}

// List returns sales with optional filters.
func (s *SaleService) List(ctx context.Context, clientID *string, date *string) ([]domain.Sale, error) {
	return s.saleRepo.FindAll(ctx, clientID, date)
}

// Create registers a new sale: calculates totals, creates credit if needed, updates client balance.
func (s *SaleService) Create(ctx context.Context, req domain.CreateSaleRequest) (*domain.Sale, error) {
	// Verify client exists
	var client *domain.Client
	var err error

	if req.ClientID == "anonymous" {
		client, err = s.clientRepo.FindByID(ctx, "anonymous")
		if err != nil {
			return nil, err
		}
		if client == nil {
			// Auto-create anonymous client
			client = &domain.Client{
				ID:          "anonymous",
				Name:        "Client de passage",
				Phone:       "",
				TotalCredit: 0,
				CreatedAt:   time.Now(),
			}
			if err := s.clientRepo.Create(ctx, client); err != nil {
				return nil, err
			}
		}
	} else {
		client, err = s.clientRepo.FindByID(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		if client == nil {
			return nil, errors.New("client not found")
		}
	}

	// Build sale items, calculate total
	var items []domain.SaleItem
	var total float64

	for _, ri := range req.Items {
		product, err := s.productRepo.FindByID(ctx, ri.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, errors.New("product not found: " + ri.ProductID)
		}

		subtotal := product.PricePerKg * ri.Quantity
		items = append(items, domain.SaleItem{
			ID:          uuid.New().String(),
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    ri.Quantity,
			Subtotal:    subtotal,
		})
		total += subtotal
	}

	// Validate paid amount
	if req.PaidAmount > total {
		return nil, errors.New("paid amount exceeds total")
	}

	creditAmount := total - req.PaidAmount

	sale := &domain.Sale{
		ID:           uuid.New().String(),
		ClientID:     client.ID,
		ClientName:   client.Name,
		Items:        items,
		Total:        total,
		PaidAmount:   req.PaidAmount,
		CreditAmount: creditAmount,
		Date:         time.Now(),
	}

	// Persist sale
	if err := s.saleRepo.Create(ctx, sale); err != nil {
		return nil, err
	}

	// If there's credit, create a credit record and update client balance
	if creditAmount > 0 {
		credit := &domain.Credit{
			ID:              uuid.New().String(),
			ClientID:        client.ID,
			ClientName:      client.Name,
			SaleID:          sale.ID,
			Amount:          creditAmount,
			RemainingAmount: creditAmount,
			Status:          domain.CreditStatusEnCours,
			CreatedAt:       time.Now(),
			Payments:        []domain.Payment{},
		}
		if err := s.creditRepo.Create(ctx, credit); err != nil {
			return nil, err
		}
		if err := s.clientRepo.UpdateTotalCredit(ctx, client.ID, creditAmount); err != nil {
			return nil, err
		}
	}

	return sale, nil
}
