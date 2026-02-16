package handler

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// SaleHandler handles HTTP requests for sale operations.
type SaleHandler struct {
	svc      *service.SaleService
	validate *validator.Validate
}

// NewSaleHandler creates a new sale handler.
func NewSaleHandler(svc *service.SaleService) *SaleHandler {
	return &SaleHandler{svc: svc, validate: validator.New()}
}

// Routes registers sale routes.
func (h *SaleHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	return r
}

func (h *SaleHandler) list(w http.ResponseWriter, r *http.Request) {
	var clientID *string
	var date *string
	if c := r.URL.Query().Get("clientId"); c != "" {
		clientID = &c
	}
	if d := r.URL.Query().Get("date"); d != "" {
		date = &d
	}
	sales, err := h.svc.List(r.Context(), clientID, date)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if sales == nil {
		sales = []domain.Sale{}
	}
	JSON(w, http.StatusOK, sales)
}

func (h *SaleHandler) create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateSaleRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	sale, err := h.svc.Create(r.Context(), req)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusCreated, sale)
}
