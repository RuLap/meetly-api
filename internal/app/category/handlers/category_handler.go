package handlers

import (
	"encoding/json"
	"github.com/RuLap/meetly-api/meetly/internal/app/category/services"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	result, err := h.categoryService.GetAllCategories(r.Context())
	if err != nil {

	}

	h.sendJSON(w, result, http.StatusOK)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {

	}

	result, err := h.categoryService.GetCategoryByID(r.Context(), id)
	if err != nil {

	}

	h.sendJSON(w, result, http.StatusOK)
}

func (h *CategoryHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
