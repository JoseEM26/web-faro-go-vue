package service_test

import (
	"errors"
	"testing"

	"proyecto-go/internal/domain"
	"proyecto-go/internal/repository"
	"proyecto-go/internal/service"
)

// mockUserRepo implementa repository.UserRepository en memoria.
type mockUserRepo struct {
	users    map[string]domain.User
	nextID   uint
	forceErr error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]domain.User), nextID: 1}
}

func (m *mockUserRepo) Create(user *domain.User) error {
	if m.forceErr != nil {
		return m.forceErr
	}
	if _, exists := m.users[user.Email]; exists {
		return repository.ErrEmailAlreadyExists
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.Email] = *user
	return nil
}

func (m *mockUserRepo) FindByEmail(email string) (domain.User, error) {
	if m.forceErr != nil {
		return domain.User{}, m.forceErr
	}
	u, ok := m.users[email]
	if !ok {
		return domain.User{}, repository.ErrNotFound
	}
	return u, nil
}

// --- Tests de Register ---

func TestRegister_Exitoso(t *testing.T) {
	svc := service.NewAuthService(newMockUserRepo(), "secreto")

	res, err := svc.Register(domain.RegisterRequest{
		Email:    "usuario@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.Token == "" {
		t.Error("se esperaba un token en la respuesta")
	}
	if res.User.Email != "usuario@test.com" {
		t.Errorf("email esperado %q, obtenido %q", "usuario@test.com", res.User.Email)
	}
	if res.User.ID == 0 {
		t.Error("se esperaba un ID de usuario generado")
	}
}

func TestRegister_EmailVacio(t *testing.T) {
	svc := service.NewAuthService(newMockUserRepo(), "secreto")

	_, err := svc.Register(domain.RegisterRequest{Email: "", Password: "password123"})
	if !errors.Is(err, service.ErrInvalidEmail) {
		t.Errorf("esperaba ErrInvalidEmail, obtuvo %v", err)
	}
}

func TestRegister_ContrasenaCorta(t *testing.T) {
	svc := service.NewAuthService(newMockUserRepo(), "secreto")

	_, err := svc.Register(domain.RegisterRequest{Email: "a@b.com", Password: "corta"})
	if !errors.Is(err, service.ErrWeakPassword) {
		t.Errorf("esperaba ErrWeakPassword, obtuvo %v", err)
	}
}

func TestRegister_EmailDuplicado(t *testing.T) {
	svc := service.NewAuthService(newMockUserRepo(), "secreto")

	svc.Register(domain.RegisterRequest{Email: "mismo@email.com", Password: "password123"})

	_, err := svc.Register(domain.RegisterRequest{Email: "mismo@email.com", Password: "otrapassword"})
	if !errors.Is(err, repository.ErrEmailAlreadyExists) {
		t.Errorf("esperaba ErrEmailAlreadyExists, obtuvo %v", err)
	}
}

// --- Tests de Login ---

func TestLogin_Exitoso(t *testing.T) {
	svc := service.NewAuthService(newMockUserRepo(), "secreto")

	// Primero registrar el usuario (el service guarda el hash de la contraseña)
	svc.Register(domain.RegisterRequest{Email: "login@test.com", Password: "mipassword"})

	res, err := svc.Login(domain.LoginRequest{
		Email:    "login@test.com",
		Password: "mipassword",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.Token == "" {
		t.Error("se esperaba un token en la respuesta")
	}
}

func TestLogin_EmailNoExiste(t *testing.T) {
	svc := service.NewAuthService(newMockUserRepo(), "secreto")

	_, err := svc.Login(domain.LoginRequest{Email: "noexiste@test.com", Password: "password"})
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Errorf("esperaba ErrInvalidCredentials, obtuvo %v", err)
	}
}

func TestLogin_ContrasenaIncorrecta(t *testing.T) {
	svc := service.NewAuthService(newMockUserRepo(), "secreto")

	svc.Register(domain.RegisterRequest{Email: "user@test.com", Password: "correcta123"})

	_, err := svc.Login(domain.LoginRequest{Email: "user@test.com", Password: "incorrecta"})
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Errorf("esperaba ErrInvalidCredentials, obtuvo %v", err)
	}
}

func TestLogin_NoRevealEmailExistence(t *testing.T) {
	svc := service.NewAuthService(newMockUserRepo(), "secreto")

	// Email que no existe debe dar el mismo error que contraseña incorrecta
	_, errEmailFaltante := svc.Login(domain.LoginRequest{Email: "fantasma@test.com", Password: "password"})
	svc.Register(domain.RegisterRequest{Email: "real@test.com", Password: "password123"})
	_, errPassFalta := svc.Login(domain.LoginRequest{Email: "real@test.com", Password: "incorrecta"})

	if errEmailFaltante.Error() != errPassFalta.Error() {
		t.Error("los errores de email no existente y password incorrecta deben ser identicos")
	}
}
