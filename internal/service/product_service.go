package service

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/port"
	"context"
	"errors"

	"github.com/google/uuid"
)

// ProductService handles product business logic.
type ProductService struct {
	repo port.ProductRepository
}

// NewProductService creates a new product service.
func NewProductService(repo port.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// List returns all products, optionally filtered by category.
func (s *ProductService) List(ctx context.Context, category *domain.MeatCategory) ([]domain.Product, error) {
	return s.repo.FindAll(ctx, category)
}

// Get returns a single product by ID.
func (s *ProductService) Get(ctx context.Context, id string) (*domain.Product, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("product not found")
	}
	return p, nil
}

// Create validates and creates a new product.
func (s *ProductService) Create(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		ID:         uuid.New().String(),
		Name:       req.Name,
		Category:   req.Category,
		PricePerKg: req.PricePerKg,
		Image:      req.Image,
		InStock:    true,
	}
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

// Update modifies an existing product's fields.
func (s *ProductService) Update(ctx context.Context, id string, req domain.UpdateProductRequest) (*domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.PricePerKg != nil {
		product.PricePerKg = *req.PricePerKg
	}
	if req.Image != nil {
		product.Image = *req.Image
	}
	if req.InStock != nil {
		product.InStock = *req.InStock
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

// Delete removes a product by ID.
func (s *ProductService) Delete(ctx context.Context, id string) error {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("product not found")
	}
	return s.repo.Delete(ctx, id)
}
