package service

import (
	"errors"

	"proyecto-go/internal/domain"
	"proyecto-go/internal/repository"
)

var ErrCategoryNameEmpty = errors.New("el nombre de la categoria no puede estar vacio")

type CategoryService interface {
	GetAll() ([]domain.Category, error)
	Create(req domain.CreateCategoryRequest) (*domain.Category, error)
	Update(id uint, req domain.UpdateCategoryRequest) (*domain.Category, error)
	Delete(id uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll() ([]domain.Category, error) {
	return s.repo.FindAll()
}

func (s *categoryService) Create(req domain.CreateCategoryRequest) (*domain.Category, error) {
	if req.Name == "" {
		return nil, ErrCategoryNameEmpty
	}
	color := req.Color
	if color == "" {
		color = "#007A6E"
	}
	cat := &domain.Category{Name: req.Name, Color: color}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Update(id uint, req domain.UpdateCategoryRequest) (*domain.Category, error) {
	if req.Name == "" {
		return nil, ErrCategoryNameEmpty
	}
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	cat.Name = req.Name
	if req.Color != "" {
		cat.Color = req.Color
	}
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}
