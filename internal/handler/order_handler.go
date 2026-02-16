package handler

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// OrderHandler handles HTTP requests for order operations.
type OrderHandler struct {
	svc      *service.OrderService
	validate *validator.Validate
}

// NewOrderHandler creates a new order handler.
func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc, validate: validator.New()}
}

// Routes registers order routes.
func (h *OrderHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	return r
}

func (h *OrderHandler) list(w http.ResponseWriter, r *http.Request) {
	var status *domain.OrderStatus
	if s := r.URL.Query().Get("status"); s != "" {
		os := domain.OrderStatus(s)
		status = &os
	}
	orders, err := h.svc.List(r.Context(), status)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if orders == nil {
		orders = []domain.Order{}
	}
	JSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	order, err := h.svc.Get(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, order)
}

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateOrderRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	order, err := h.svc.Create(r.Context(), req)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateOrderRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	order, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, order)
}

func (h *OrderHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"deleted": id})
}
