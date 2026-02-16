package service

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/port"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// OrderService handles order business logic.
type OrderService struct {
	orderRepo   port.OrderRepository
	clientRepo  port.ClientRepository
	productRepo port.ProductRepository
}

// NewOrderService creates a new order service.
func NewOrderService(orderRepo port.OrderRepository, clientRepo port.ClientRepository, productRepo port.ProductRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, clientRepo: clientRepo, productRepo: productRepo}
}

// List returns orders, optionally filtered by status.
func (s *OrderService) List(ctx context.Context, status *domain.OrderStatus) ([]domain.Order, error) {
	return s.orderRepo.FindAll(ctx, status)
}

// Get returns a single order by ID.
func (s *OrderService) Get(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

// Create validates and creates a new order.
func (s *OrderService) Create(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error) {
	client, err := s.clientRepo.FindByID(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client not found")
	}

	pickupDate, err := time.Parse("2006-01-02", req.PickupDate)
	if err != nil {
		return nil, errors.New("invalid pickup date format, expected YYYY-MM-DD")
	}

	var items []domain.OrderItem
	for _, ri := range req.Items {
		product, err := s.productRepo.FindByID(ctx, ri.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, errors.New("product not found: " + ri.ProductID)
		}
		items = append(items, domain.OrderItem{
			ID:          uuid.New().String(),
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    ri.Quantity,
		})
	}

	order := &domain.Order{
		ID:          uuid.New().String(),
		ClientID:    client.ID,
		ClientName:  client.Name,
		ClientPhone: client.Phone,
		Items:       items,
		PickupDate:  pickupDate,
		Notes:       req.Notes,
		Status:      domain.OrderStatusEnAttente,
		CreatedAt:   time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

// Update modifies an existing order's status, date, or notes.
func (s *OrderService) Update(ctx context.Context, id string, req domain.UpdateOrderRequest) (*domain.Order, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if req.Status != nil {
		order.Status = *req.Status
	}
	if req.PickupDate != nil {
		d, err := time.Parse("2006-01-02", *req.PickupDate)
		if err != nil {
			return nil, errors.New("invalid pickup date format")
		}
		order.PickupDate = d
	}
	if req.Notes != nil {
		order.Notes = *req.Notes
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

// Delete removes an order by ID.
func (s *OrderService) Delete(ctx context.Context, id string) error {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}
	return s.orderRepo.Delete(ctx, id)
}
