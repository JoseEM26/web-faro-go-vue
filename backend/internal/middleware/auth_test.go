package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"proyecto-go/internal/middleware"
	"proyecto-go/pkg/token"
)

const testSecret = "secreto-de-test"

func handlerOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuth_SinHeader(t *testing.T) {
	handler := middleware.Auth(testSecret)(http.HandlerFunc(handlerOK))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuvo %d", rec.Code)
	}
}

func TestAuth_HeaderSinBearer(t *testing.T) {
	handler := middleware.Auth(testSecret)(http.HandlerFunc(handlerOK))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuvo %d", rec.Code)
	}
}

func TestAuth_TokenValido(t *testing.T) {
	tokenStr, err := token.Generate(7, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("error generando token: %v", err)
	}

	handler := middleware.Auth(testSecret)(http.HandlerFunc(handlerOK))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperaba 200, obtuvo %d", rec.Code)
	}
}

func TestAuth_TokenInvalido(t *testing.T) {
	handler := middleware.Auth(testSecret)(http.HandlerFunc(handlerOK))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer token.completamente.invalido")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuvo %d", rec.Code)
	}
}

func TestAuth_TokenSecretoDistinto(t *testing.T) {
	tokenStr, _ := token.Generate(1, "secreto-A", time.Hour)

	handler := middleware.Auth("secreto-B")(http.HandlerFunc(handlerOK))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuvo %d", rec.Code)
	}
}

func TestAuth_TokenExpirado(t *testing.T) {
	tokenStr, _ := token.Generate(1, testSecret, -time.Hour)

	handler := middleware.Auth(testSecret)(http.HandlerFunc(handlerOK))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuvo %d", rec.Code)
	}
}

func TestAuth_UserIDEnContexto(t *testing.T) {
	tokenStr, _ := token.Generate(99, testSecret, time.Hour)

	var capturedID uint
	handler := middleware.Auth(testSecret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = r.Context().Value(middleware.UserIDKey).(uint)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if capturedID != 99 {
		t.Errorf("UserID esperado 99, obtenido %d", capturedID)
	}
}
