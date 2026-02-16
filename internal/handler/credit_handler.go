package handler

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// CreditHandler handles HTTP requests for credit operations.
type CreditHandler struct {
	svc      *service.CreditService
	validate *validator.Validate
}

// NewCreditHandler creates a new credit handler.
func NewCreditHandler(svc *service.CreditService) *CreditHandler {
	return &CreditHandler{svc: svc, validate: validator.New()}
}

// Routes registers credit routes.
func (h *CreditHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/{id}/payments", h.addPayment)
	return r
}

func (h *CreditHandler) list(w http.ResponseWriter, r *http.Request) {
	var status *domain.CreditStatus
	if s := r.URL.Query().Get("status"); s != "" {
		cs := domain.CreditStatus(s)
		status = &cs
	}
	credits, err := h.svc.List(r.Context(), status)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if credits == nil {
		credits = []domain.Credit{}
	}
	JSON(w, http.StatusOK, credits)
}

func (h *CreditHandler) addPayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.CreatePaymentRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	credit, err := h.svc.AddPayment(r.Context(), id, req)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusOK, credit)
}
