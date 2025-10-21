package user

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"

	validation "github.com/RuLap/meetly-api/meetly/internal/pkg/validator"
	"github.com/darahayes/go-boom"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		boom.BadRequest(w, "ID обязателен")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		boom.BadRequest(w, "неверный формат ID")
		return
	}

	response, err := h.service.GetByID(r.Context(), uid)
	if err != nil {
		boom.Internal(w, err.Error())
		return
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		boom.BadRequest(w, "ID обязателен")
		return
	}

	var req SaveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		boom.BadRequest(w, err.Error())
		return
	}

	if errors := validation.ValidateStruct(req); errors != nil {
		boom.BadRequest(w, "Ошибки валидации", errors)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		boom.BadRequest(w, "неверный формат ID")
		return
	}

	response, err := h.service.UpdateUser(r.Context(), uid, &req)
	if err != nil {
		boom.Internal(w, err.Error())
		return
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *Handler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
