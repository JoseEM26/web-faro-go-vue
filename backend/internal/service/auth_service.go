package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"proyecto-go/internal/domain"
	"proyecto-go/internal/repository"
	"proyecto-go/pkg/token"
)

var (
	ErrInvalidEmail       = errors.New("el email no puede estar vacio")
	ErrWeakPassword       = errors.New("la contrasena debe tener al menos 8 caracteres")
	ErrInvalidCredentials = errors.New("email o contrasena incorrectos")
)

type AuthService interface {
	Register(req domain.RegisterRequest) (domain.AuthResponse, error)
	Login(req domain.LoginRequest) (domain.AuthResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtTTL:    24 * time.Hour,
	}
}

func (s *authService) Register(req domain.RegisterRequest) (domain.AuthResponse, error) {
	if req.Email == "" {
		return domain.AuthResponse{}, ErrInvalidEmail
	}
	if len(req.Password) < 8 {
		return domain.AuthResponse{}, ErrWeakPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	user := domain.User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}
	if err := s.userRepo.Create(&user); err != nil {
		return domain.AuthResponse{}, err
	}

	t, err := token.Generate(user.ID, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	return domain.AuthResponse{
		Token: t,
		User:  domain.UserDTO{ID: user.ID, Email: user.Email},
	}, nil
}

func (s *authService) Login(req domain.LoginRequest) (domain.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Se devuelve siempre el mismo error para no revelar si el email existe
		return domain.AuthResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return domain.AuthResponse{}, ErrInvalidCredentials
	}

	t, err := token.Generate(user.ID, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	return domain.AuthResponse{
		Token: t,
		User:  domain.UserDTO{ID: user.ID, Email: user.Email},
	}, nil
}
