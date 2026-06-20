package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"proyecto-go/internal/domain"
)

var ErrEmailAlreadyExists = errors.New("el email ya esta registrado")

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (domain.User, error)
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if isUniqueViolation(err) {
			return ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *gormUserRepository) FindByEmail(email string) (domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}

// isUniqueViolation detecta violaciones de restriccion UNIQUE de PostgreSQL (codigo 23505).
func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key")
}
