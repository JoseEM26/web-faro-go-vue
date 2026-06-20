package service_test

import (
	"errors"
	"testing"

	"proyecto-go/internal/domain"
	"proyecto-go/internal/repository"
	"proyecto-go/internal/service"
)

// mockRepo implementa repository.TaskRepository en memoria.
// Permite testear el service sin base de datos.
type mockRepo struct {
	tasks    map[uint]domain.Task
	nextID   uint
	forceErr error
}

func newMockRepo() *mockRepo {
	return &mockRepo{tasks: make(map[uint]domain.Task), nextID: 1}
}

func (m *mockRepo) FindAll() ([]domain.Task, error) {
	if m.forceErr != nil {
		return nil, m.forceErr
	}
	tasks := make([]domain.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (m *mockRepo) FindByID(id uint) (domain.Task, error) {
	if m.forceErr != nil {
		return domain.Task{}, m.forceErr
	}
	t, ok := m.tasks[id]
	if !ok {
		return domain.Task{}, repository.ErrNotFound
	}
	return t, nil
}

func (m *mockRepo) Create(task *domain.Task) error {
	if m.forceErr != nil {
		return m.forceErr
	}
	task.ID = m.nextID
	m.nextID++
	m.tasks[task.ID] = *task
	return nil
}

func (m *mockRepo) Update(task *domain.Task) error {
	if m.forceErr != nil {
		return m.forceErr
	}
	if _, ok := m.tasks[task.ID]; !ok {
		return repository.ErrNotFound
	}
	m.tasks[task.ID] = *task
	return nil
}

func (m *mockRepo) Delete(id uint) error {
	if m.forceErr != nil {
		return m.forceErr
	}
	if _, ok := m.tasks[id]; !ok {
		return repository.ErrNotFound
	}
	delete(m.tasks, id)
	return nil
}

// --- Tests ---

func setup() service.TaskService {
	return service.NewTaskService(newMockRepo())
}

func TestCreate_Valido(t *testing.T) {
	svc := setup()

	task, err := svc.Create(domain.CreateTaskRequest{Title: "Mi tarea", Description: "desc"})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if task.Title != "Mi tarea" {
		t.Errorf("titulo esperado %q, obtenido %q", "Mi tarea", task.Title)
	}
	if task.ID == 0 {
		t.Error("se esperaba un ID generado")
	}
	if task.Completed {
		t.Error("nueva tarea no debe estar completada")
	}
}

func TestCreate_TituloVacio(t *testing.T) {
	svc := setup()

	_, err := svc.Create(domain.CreateTaskRequest{Title: ""})
	if !errors.Is(err, service.ErrInvalidTitle) {
		t.Errorf("esperaba ErrInvalidTitle, obtuvo %v", err)
	}
}

func TestGetAll_Vacio(t *testing.T) {
	svc := setup()

	tasks, err := svc.GetAll()
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("esperaba 0 tareas, obtuvo %d", len(tasks))
	}
}

func TestGetAll_ConTareas(t *testing.T) {
	svc := setup()
	svc.Create(domain.CreateTaskRequest{Title: "Tarea 1"})
	svc.Create(domain.CreateTaskRequest{Title: "Tarea 2"})

	tasks, err := svc.GetAll()
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("esperaba 2 tareas, obtuvo %d", len(tasks))
	}
}

func TestGetByID_Exitoso(t *testing.T) {
	svc := setup()
	created, _ := svc.Create(domain.CreateTaskRequest{Title: "Tarea"})

	found, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if found.ID != created.ID {
		t.Errorf("ID esperado %d, obtenido %d", created.ID, found.ID)
	}
}

func TestGetByID_NoEncontrado(t *testing.T) {
	svc := setup()

	_, err := svc.GetByID(999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("esperaba ErrNotFound, obtuvo %v", err)
	}
}

func TestUpdate_Exitoso(t *testing.T) {
	svc := setup()
	created, _ := svc.Create(domain.CreateTaskRequest{Title: "Original"})

	updated, err := svc.Update(created.ID, domain.UpdateTaskRequest{
		Title:     "Actualizada",
		Completed: true,
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if updated.Title != "Actualizada" {
		t.Errorf("titulo esperado %q, obtenido %q", "Actualizada", updated.Title)
	}
	if !updated.Completed {
		t.Error("se esperaba que la tarea estuviera completada")
	}
}

func TestUpdate_NoEncontrado(t *testing.T) {
	svc := setup()

	_, err := svc.Update(999, domain.UpdateTaskRequest{Title: "Fantasma"})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("esperaba ErrNotFound, obtuvo %v", err)
	}
}

func TestDelete_Exitoso(t *testing.T) {
	svc := setup()
	created, _ := svc.Create(domain.CreateTaskRequest{Title: "A eliminar"})

	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	_, err := svc.GetByID(created.ID)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Error("se esperaba ErrNotFound al buscar tarea eliminada")
	}
}

func TestDelete_NoEncontrado(t *testing.T) {
	svc := setup()

	err := svc.Delete(999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("esperaba ErrNotFound, obtuvo %v", err)
	}
}
