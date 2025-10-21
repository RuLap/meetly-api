package event

import (
	"encoding/json"
	"net/http"

	validation "github.com/RuLap/meetly-api/meetly/internal/pkg/validator"
	"github.com/darahayes/go-boom"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetShortEvents(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetShortEvents(r.Context())
	if err != nil {
		boom.Internal(w, err)
		return
	}

	h.sendJSON(w, result, http.StatusOK)
}

func (h *Handler) GetEventWithDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		boom.BadRequest(w, "ID необходим")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		boom.BadRequest(w, "неверный формат ID")
		return
	}

	result, err := h.service.GetEventWithDetails(r.Context(), uid)
	if err != nil {
		boom.Internal(w, err)
		return
	}

	h.sendJSON(w, result, http.StatusOK)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		boom.BadRequest(w, "неверный формат JSON")
		return
	}

	if errors := validation.ValidateStruct(req); errors != nil {
		boom.BadRequest(w, "Ошибки валидации", errors)
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	result, err := h.service.CreateEvent(r.Context(), &req, userID)
	if err != nil {
		boom.Internal(w, err)
		return
	}

	h.sendJSON(w, result, http.StatusOK)
}

func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetAllCategories(r.Context())
	if err != nil {
		boom.Internal(w, err)
		return
	}

	h.sendJSON(w, result, http.StatusOK)
}

func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		boom.BadRequest(w, "ID необходим")
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		boom.BadRequest(w, "неверный формат ID")
		return
	}

	result, err := h.service.GetCategoryByID(r.Context(), uid)
	if err != nil {
		boom.Internal(w, err)
		return
	}

	h.sendJSON(w, result, http.StatusOK)
}

func (h *Handler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	if eventID == "" {
		boom.BadRequest(w, "ID необходим")
		return
	}
	uid, err := uuid.Parse(eventID)
	if err != nil {
		boom.BadRequest(w, "неверный формат ID")
		return
	}

	userID := r.Context().Value("user_id").(uuid.UUID)

	result, err := h.service.AddParticipant(r.Context(), uid, userID)
	if err != nil {
		boom.Internal(w, err)
		return
	}

	h.sendJSON(w, result, http.StatusOK)
}

func (h *Handler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
