package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"proyecto-go/internal/domain"
	"proyecto-go/internal/handler"
	"proyecto-go/internal/repository"
	"proyecto-go/internal/service"
)

// mockAuthService implementa service.AuthService con respuestas controladas.
type mockAuthService struct {
	response  domain.AuthResponse
	returnErr error
}

func (m *mockAuthService) Register(req domain.RegisterRequest) (domain.AuthResponse, error) {
	return m.response, m.returnErr
}
func (m *mockAuthService) Login(req domain.LoginRequest) (domain.AuthResponse, error) {
	return m.response, m.returnErr
}

func setupAuthRouter(svc *mockAuthService) *chi.Mux {
	r := chi.NewRouter()
	h := handler.NewAuthHandler(svc)
	r.Route("/auth", h.RegisterRoutes)
	return r
}

func TestRegister_Exitoso(t *testing.T) {
	svc := &mockAuthService{
		response: domain.AuthResponse{
			Token: "un.jwt.token",
			User:  domain.UserDTO{ID: 1, Email: "nuevo@test.com"},
		},
	}
	r := setupAuthRouter(svc)

	body, _ := json.Marshal(domain.RegisterRequest{Email: "nuevo@test.com", Password: "password123"})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("esperaba status 201, obtuvo %d", rec.Code)
	}

	var res domain.AuthResponse
	json.NewDecoder(rec.Body).Decode(&res)
	if res.Token == "" {
		t.Error("se esperaba un token en la respuesta")
	}
}

func TestRegister_CuerpoInvalido(t *testing.T) {
	r := setupAuthRouter(&mockAuthService{})

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("no es json")))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperaba status 400, obtuvo %d", rec.Code)
	}
}

func TestRegister_EmailVacio(t *testing.T) {
	svc := &mockAuthService{returnErr: service.ErrInvalidEmail}
	r := setupAuthRouter(svc)

	body, _ := json.Marshal(domain.RegisterRequest{Email: "", Password: "password123"})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperaba status 400, obtuvo %d", rec.Code)
	}
}

func TestRegister_ContrasenaCorta(t *testing.T) {
	svc := &mockAuthService{returnErr: service.ErrWeakPassword}
	r := setupAuthRouter(svc)

	body, _ := json.Marshal(domain.RegisterRequest{Email: "a@b.com", Password: "corta"})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperaba status 400, obtuvo %d", rec.Code)
	}
}

func TestRegister_EmailDuplicado(t *testing.T) {
	svc := &mockAuthService{returnErr: repository.ErrEmailAlreadyExists}
	r := setupAuthRouter(svc)

	body, _ := json.Marshal(domain.RegisterRequest{Email: "duplicado@test.com", Password: "password123"})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("esperaba status 409, obtuvo %d", rec.Code)
	}
}

func TestLogin_Exitoso(t *testing.T) {
	svc := &mockAuthService{
		response: domain.AuthResponse{
			Token: "un.jwt.token",
			User:  domain.UserDTO{ID: 1, Email: "user@test.com"},
		},
	}
	r := setupAuthRouter(svc)

	body, _ := json.Marshal(domain.LoginRequest{Email: "user@test.com", Password: "password123"})
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperaba status 200, obtuvo %d", rec.Code)
	}

	var res domain.AuthResponse
	json.NewDecoder(rec.Body).Decode(&res)
	if res.Token == "" {
		t.Error("se esperaba un token en la respuesta")
	}
}

func TestLogin_CredencialesInvalidas(t *testing.T) {
	svc := &mockAuthService{returnErr: service.ErrInvalidCredentials}
	r := setupAuthRouter(svc)

	body, _ := json.Marshal(domain.LoginRequest{Email: "user@test.com", Password: "incorrecta"})
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperaba status 401, obtuvo %d", rec.Code)
	}
}

func TestLogin_CuerpoInvalido(t *testing.T) {
	r := setupAuthRouter(&mockAuthService{})

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("no es json")))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperaba status 400, obtuvo %d", rec.Code)
	}
}
