package service

import (
	"errors"

	"proyecto-go/internal/domain"
	"proyecto-go/internal/repository"
)

var ErrNoteTitleEmpty = errors.New("el titulo de la nota no puede estar vacio")

type NoteService interface {
	GetAll(userID uint) ([]domain.Note, error)
	GetByID(id, userID uint) (*domain.Note, error)
	Create(userID uint, req domain.CreateNoteRequest) (*domain.Note, error)
	Update(id, userID uint, req domain.UpdateNoteRequest) (*domain.Note, error)
	Delete(id, userID uint) error
}

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) GetAll(userID uint) ([]domain.Note, error) {
	return s.repo.FindAll(userID)
}

func (s *noteService) GetByID(id, userID uint) (*domain.Note, error) {
	return s.repo.FindByID(id, userID)
}

func (s *noteService) Create(userID uint, req domain.CreateNoteRequest) (*domain.Note, error) {
	if req.Title == "" {
		return nil, ErrNoteTitleEmpty
	}
	note := &domain.Note{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}
	if err := s.repo.Create(note); err != nil {
		return nil, err
	}
	return note, nil
}

func (s *noteService) Update(id, userID uint, req domain.UpdateNoteRequest) (*domain.Note, error) {
	if req.Title == "" {
		return nil, ErrNoteTitleEmpty
	}
	note, err := s.repo.FindByID(id, userID)
	if err != nil {
		return nil, err
	}
	note.Title = req.Title
	note.Content = req.Content
	if err := s.repo.Update(note); err != nil {
		return nil, err
	}
	return note, nil
}

func (s *noteService) Delete(id, userID uint) error {
	return s.repo.Delete(id, userID)
}
