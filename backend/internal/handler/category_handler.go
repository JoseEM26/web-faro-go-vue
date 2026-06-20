package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"proyecto-go/internal/domain"
	"proyecto-go/internal/repository"
	"proyecto-go/internal/service"
	"proyecto-go/pkg/response"
)

type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.GetAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "error obteniendo categorias")
		return
	}
	response.JSON(w, http.StatusOK, cats)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "cuerpo de solicitud invalido")
		return
	}
	cat, err := h.svc.Create(req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNameEmpty) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "error creando categoria")
		return
	}
	response.JSON(w, http.StatusCreated, cat)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "id invalido")
		return
	}
	var req domain.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "cuerpo de solicitud invalido")
		return
	}
	cat, err := h.svc.Update(id, req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNameEmpty) {
			response.Error(w, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "categoria no encontrada")
		} else {
			response.Error(w, http.StatusInternalServerError, "error actualizando categoria")
		}
		return
	}
	response.JSON(w, http.StatusOK, cat)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "id invalido")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.Error(w, http.StatusInternalServerError, "error eliminando categoria")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
