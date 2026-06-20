package service

import (
	"errors"

	"proyecto-go/internal/domain"
	"proyecto-go/internal/repository"
)

var ErrInvalidTitle = errors.New("el titulo no puede estar vacio")

type TaskService interface {
	GetAll() ([]domain.Task, error)
	GetByID(id uint) (domain.Task, error)
	Create(req domain.CreateTaskRequest) (domain.Task, error)
	Update(id uint, req domain.UpdateTaskRequest) (domain.Task, error)
	Delete(id uint) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) GetAll() ([]domain.Task, error) {
	return s.repo.FindAll()
}

func (s *taskService) GetByID(id uint) (domain.Task, error) {
	return s.repo.FindByID(id)
}

func (s *taskService) Create(req domain.CreateTaskRequest) (domain.Task, error) {
	if req.Title == "" {
		return domain.Task{}, ErrInvalidTitle
	}

	task := domain.Task{
		Title:       req.Title,
		Description: req.Description,
	}
	if err := s.repo.Create(&task); err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (s *taskService) Update(id uint, req domain.UpdateTaskRequest) (domain.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Task{}, err
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Completed = req.Completed

	if err := s.repo.Update(&task); err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (s *taskService) Delete(id uint) error {
	return s.repo.Delete(id)
}
