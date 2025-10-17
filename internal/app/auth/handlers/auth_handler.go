package handlers

import (
	"encoding/json"
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/dto"
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/services"
	validation "github.com/RuLap/meetly-api/meetly/internal/pkg/validator"
	"github.com/darahayes/go-boom"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		boom.BadRequest(w, "неверный формат JSON")
		return
	}

	if errors := validation.ValidateStruct(req); errors != nil {
		boom.BadRequest(w, "Ошибки валидации", errors)
		return
	}

	response, err := h.authService.Register(r.Context(), req)
	if err != nil {
		boom.BadRequest(w, err.Error())
		return
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		boom.BadRequest(w, err.Error())
		return
	}

	if errors := validation.ValidateStruct(req); errors != nil {
		boom.BadRequest(w, "Ошибки валидации", errors)
		return
	}

	response, err := h.authService.Login(r.Context(), req)
	if err != nil {
		boom.BadRequest(w, err.Error())
		return
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *AuthHandler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req dto.GoogleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		boom.BadRequest(w, "Неверный формат JSON")
		return
	}

	if errors := validation.ValidateStruct(req); errors != nil {
		boom.BadRequest(w, "Ошибки валидации", errors)
		return
	}

	response, err := h.authService.GoogleAuth(r.Context(), req)
	if err != nil {
		boom.Unathorized(w, err.Error())
		return
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *AuthHandler) GoogleAuthURL(w http.ResponseWriter, r *http.Request) {
	url, state, err := h.authService.GenerateGoogleOAuthURL()
	if err != nil {
		boom.Internal(w, "Не удалось сгенерировать URL для авторизации")
		return
	}

	response := map[string]string{
		"url":   url,
		"state": state,
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *AuthHandler) SendConfirmationLink(w http.ResponseWriter, r *http.Request) {
	var req dto.ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		boom.BadRequest(w, "неверный формат JSON")
		return
	}

	if errors := validation.ValidateStruct(req); errors != nil {
		boom.BadRequest(w, "Ошибки валидации", errors)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		boom.Unathorized(w, "требуется аутентификация")
		return
	}

	if err := h.authService.ConfirmEmail(r.Context(), req.Token, userID); err != nil {
		boom.BadRequest(w, err.Error())
		return
	}

	h.sendJSON(w, map[string]interface{}{
		"success": true,
		"message": "Email успешно подтвержден",
	}, http.StatusOK)
}

func (h *AuthHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	var req dto.ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		boom.BadRequest(w, "неверный формат JSON")
		return
	}

	if errors := validation.ValidateStruct(req); errors != nil {
		boom.BadRequest(w, "Ошибки валидации", errors)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		boom.Unathorized(w, "требуется аутентификация")
		return
	}

	if err := h.authService.ConfirmEmail(r.Context(), req.Token, userID); err != nil {
		boom.BadRequest(w, err.Error())
		return
	}

	h.sendJSON(w, map[string]interface{}{
		"success": true,
		"message": "Email успешно подтвержден",
	}, http.StatusOK)
}

func (h *AuthHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
