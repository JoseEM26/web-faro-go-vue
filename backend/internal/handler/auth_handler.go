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

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "cuerpo de solicitud invalido")
		return
	}

	res, err := h.svc.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidEmail) || errors.Is(err, service.ErrWeakPassword) {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "error interno")
		return
	}

	response.JSON(w, http.StatusCreated, res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "cuerpo de solicitud invalido")
		return
	}

	res, err := h.svc.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "error interno")
		return
	}

	response.JSON(w, http.StatusOK, res)
}
