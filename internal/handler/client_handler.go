package handler

import (
	"boucherie-api/internal/domain"
	"boucherie-api/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// ClientHandler handles HTTP requests for client operations.
type ClientHandler struct {
	svc      *service.ClientService
	validate *validator.Validate
}

// NewClientHandler creates a new client handler.
func NewClientHandler(svc *service.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc, validate: validator.New()}
}

// Routes registers client routes on the given router.
func (h *ClientHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Get("/{id}", h.get)
	r.Put("/{id}", h.update)
	return r
}

func (h *ClientHandler) list(w http.ResponseWriter, r *http.Request) {
	clients, err := h.svc.List(r.Context())
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if clients == nil {
		clients = []domain.Client{}
	}
	JSON(w, http.StatusOK, clients)
}

func (h *ClientHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	client, err := h.svc.Get(r.Context(), id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, client)
}

func (h *ClientHandler) create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateClientRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	client, err := h.svc.Create(r.Context(), req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, client)
}

func (h *ClientHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateClientRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	client, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, client)
}
