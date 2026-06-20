package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"proyecto-go/internal/domain"
	"proyecto-go/internal/middleware"
	"proyecto-go/internal/repository"
	"proyecto-go/internal/service"
	"proyecto-go/pkg/response"
)

type NoteHandler struct {
	svc service.NoteService
}

func NewNoteHandler(svc service.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

func (h *NoteHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

func userIDFromCtx(r *http.Request) uint {
	return r.Context().Value(middleware.UserIDKey).(uint)
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	notes, err := h.svc.GetAll(userIDFromCtx(r))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "error obteniendo notas")
		return
	}
	response.JSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "id invalido")
		return
	}
	note, err := h.svc.GetByID(id, userIDFromCtx(r))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "nota no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "error obteniendo nota")
		return
	}
	response.JSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "cuerpo de solicitud invalido")
		return
	}
	note, err := h.svc.Create(userIDFromCtx(r), req)
	if err != nil {
		if errors.Is(err, service.ErrNoteTitleEmpty) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "error creando nota")
		return
	}
	response.JSON(w, http.StatusCreated, note)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "id invalido")
		return
	}
	var req domain.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "cuerpo de solicitud invalido")
		return
	}
	note, err := h.svc.Update(id, userIDFromCtx(r), req)
	if err != nil {
		if errors.Is(err, service.ErrNoteTitleEmpty) {
			response.Error(w, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "nota no encontrada")
		} else {
			response.Error(w, http.StatusInternalServerError, "error actualizando nota")
		}
		return
	}
	response.JSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "id invalido")
		return
	}
	if err := h.svc.Delete(id, userIDFromCtx(r)); err != nil {
		response.Error(w, http.StatusInternalServerError, "error eliminando nota")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
