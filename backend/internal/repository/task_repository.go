package repository

import (
	"errors"

	"gorm.io/gorm"
	"proyecto-go/internal/domain"
)

var ErrNotFound = errors.New("registro no encontrado")

// TaskRepository define las operaciones de acceso a datos.
// El service solo conoce esta interfaz, nunca el *gorm.DB.
type TaskRepository interface {
	FindAll() ([]domain.Task, error)
	FindByID(id uint) (domain.Task, error)
	Create(task *domain.Task) error
	Update(task *domain.Task) error
	Delete(id uint) error
}

type gormTaskRepository struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) TaskRepository {
	return &gormTaskRepository{db: db}
}

func (r *gormTaskRepository) FindAll() ([]domain.Task, error) {
	var tasks []domain.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *gormTaskRepository) FindByID(id uint) (domain.Task, error) {
	var task domain.Task
	if err := r.db.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Task{}, ErrNotFound
		}
		return domain.Task{}, err
	}
	return task, nil
}

func (r *gormTaskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *gormTaskRepository) Update(task *domain.Task) error {
	return r.db.Save(task).Error
}

func (r *gormTaskRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
