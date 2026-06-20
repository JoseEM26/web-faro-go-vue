package repository

import (
	"errors"

	"gorm.io/gorm"

	"proyecto-go/internal/domain"
)

type NoteRepository interface {
	FindAll(userID uint) ([]domain.Note, error)
	FindByID(id, userID uint) (*domain.Note, error)
	Create(note *domain.Note) error
	Update(note *domain.Note) error
	Delete(id, userID uint) error
}

type gormNoteRepository struct {
	db *gorm.DB
}

func NewGormNoteRepository(db *gorm.DB) NoteRepository {
	return &gormNoteRepository{db: db}
}

func (r *gormNoteRepository) FindAll(userID uint) ([]domain.Note, error) {
	var notes []domain.Note
	if err := r.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *gormNoteRepository) FindByID(id, userID uint) (*domain.Note, error) {
	var note domain.Note
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (r *gormNoteRepository) Create(note *domain.Note) error {
	return r.db.Create(note).Error
}

func (r *gormNoteRepository) Update(note *domain.Note) error {
	return r.db.Save(note).Error
}

func (r *gormNoteRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Note{}).Error
}
