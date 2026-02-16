package handler

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// ProductHandler handles HTTP requests for product operations.
type ProductHandler struct {
	svc      *service.ProductService
	validate *validator.Validate
}

// NewProductHandler creates a new product handler.
func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc, validate: validator.New()}
}

// Routes registers product routes.
func (h *ProductHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
	return r
}

func (h *ProductHandler) list(w http.ResponseWriter, r *http.Request) {
	var category *domain.MeatCategory
	if c := r.URL.Query().Get("category"); c != "" {
		cat := domain.MeatCategory(c)
		category = &cat
	}
	products, err := h.svc.List(r.Context(), category)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if products == nil {
		products = []domain.Product{}
	}
	JSON(w, http.StatusOK, products)
}

func (h *ProductHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := h.svc.Get(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProductRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	product, err := h.svc.Create(r.Context(), req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateProductRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	product, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"deleted": id})
}
