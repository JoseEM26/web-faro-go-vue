package repository

import (
	"errors"

	"gorm.io/gorm"

	"proyecto-go/internal/domain"
)

type CategoryRepository interface {
	FindAll() ([]domain.Category, error)
	FindByID(id uint) (*domain.Category, error)
	Create(cat *domain.Category) error
	Update(cat *domain.Category) error
	Delete(id uint) error
}

type gormCategoryRepository struct {
	db *gorm.DB
}

func NewGormCategoryRepository(db *gorm.DB) CategoryRepository {
	return &gormCategoryRepository{db: db}
}

func (r *gormCategoryRepository) FindAll() ([]domain.Category, error) {
	var cats []domain.Category
	if err := r.db.Order("name asc").Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

func (r *gormCategoryRepository) FindByID(id uint) (*domain.Category, error) {
	var cat domain.Category
	if err := r.db.First(&cat, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &cat, nil
}

func (r *gormCategoryRepository) Create(cat *domain.Category) error {
	return r.db.Create(cat).Error
}

func (r *gormCategoryRepository) Update(cat *domain.Category) error {
	return r.db.Save(cat).Error
}

func (r *gormCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}
