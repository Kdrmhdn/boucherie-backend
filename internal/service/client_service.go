package service

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/port"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ClientService handles client business logic.
type ClientService struct {
	repo port.ClientRepository
}

// NewClientService creates a new client service.
func NewClientService(repo port.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

// List returns all clients.
func (s *ClientService) List(ctx context.Context) ([]domain.Client, error) {
	return s.repo.FindAll(ctx)
}

// Get returns a single client by ID.
func (s *ClientService) Get(ctx context.Context, id string) (*domain.Client, error) {
	client, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client not found")
	}
	return client, nil
}

// Create validates and creates a new client.
func (s *ClientService) Create(ctx context.Context, req domain.CreateClientRequest) (*domain.Client, error) {
	client := &domain.Client{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Phone:       req.Phone,
		Email:       req.Email,
		TotalCredit: 0,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}

// Update modifies an existing client's fields.
func (s *ClientService) Update(ctx context.Context, id string, req domain.UpdateClientRequest) (*domain.Client, error) {
	client, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client not found")
	}

	if req.Name != nil {
		client.Name = *req.Name
	}
	if req.Phone != nil {
		client.Phone = *req.Phone
	}
	if req.Email != nil {
		client.Email = *req.Email
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return nil, err
	}
	return client, nil
}
