package handlers

import (
	"encoding/json"
	"github.com/RuLap/meetly-api/meetly/internal/app/user/dto"
	"github.com/RuLap/meetly-api/meetly/internal/app/user/services"
	validation "github.com/RuLap/meetly-api/meetly/internal/pkg/validator"
	"github.com/darahayes/go-boom"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		boom.BadRequest(w, "ID обязателен")
		return
	}

	response, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		boom.Internal(w, err.Error())
		return
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		boom.BadRequest(w, "ID обязателен")
		return
	}

	var req dto.SaveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		boom.BadRequest(w, err.Error())
		return
	}

	if errors := validation.ValidateStruct(req); errors != nil {
		boom.BadRequest(w, "Ошибки валидации", errors)
		return
	}

	response, err := h.userService.UpdateUser(r.Context(), id, &req)
	if err != nil {
		boom.Internal(w, err.Error())
		return
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *UserHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
