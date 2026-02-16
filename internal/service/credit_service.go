package service

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/port"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// CreditService handles credit business logic.
type CreditService struct {
	creditRepo port.CreditRepository
	clientRepo port.ClientRepository
}

// NewCreditService creates a new credit service.
func NewCreditService(creditRepo port.CreditRepository, clientRepo port.ClientRepository) *CreditService {
	return &CreditService{creditRepo: creditRepo, clientRepo: clientRepo}
}

// List returns all credits, optionally filtered by status.
func (s *CreditService) List(ctx context.Context, status *domain.CreditStatus) ([]domain.Credit, error) {
	return s.creditRepo.FindAll(ctx, status)
}

// ListByClient returns credits for a specific client.
func (s *CreditService) ListByClient(ctx context.Context, clientID string) ([]domain.Credit, error) {
	return s.creditRepo.FindByClientID(ctx, clientID)
}

// AddPayment registers a payment on a credit, updates remaining amount, and adjusts client balance.
func (s *CreditService) AddPayment(ctx context.Context, creditID string, req domain.CreatePaymentRequest) (*domain.Credit, error) {
	credit, err := s.creditRepo.FindByID(ctx, creditID)
	if err != nil {
		return nil, err
	}
	if credit == nil {
		return nil, errors.New("credit not found")
	}
	if credit.Status == domain.CreditStatusPaye {
		return nil, errors.New("credit already fully paid")
	}
	if req.Amount > credit.RemainingAmount {
		return nil, errors.New("payment exceeds remaining amount")
	}

	// Create payment
	payment := &domain.Payment{
		ID:       uuid.New().String(),
		CreditID: creditID,
		Amount:   req.Amount,
		Date:     time.Now(),
		Method:   req.Method,
	}
	if err := s.creditRepo.AddPayment(ctx, payment); err != nil {
		return nil, err
	}

	// Update credit
	credit.RemainingAmount -= req.Amount
	if credit.RemainingAmount <= 0 {
		credit.RemainingAmount = 0
		credit.Status = domain.CreditStatusPaye
	}
	if err := s.creditRepo.Update(ctx, credit); err != nil {
		return nil, err
	}

	// Update client balance (decrease)
	if err := s.clientRepo.UpdateTotalCredit(ctx, credit.ClientID, -req.Amount); err != nil {
		return nil, err
	}

	// Reload with payments
	return s.creditRepo.FindByID(ctx, creditID)
}
