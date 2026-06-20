package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"proyecto-go/internal/domain"
	"proyecto-go/internal/handler"
	"proyecto-go/internal/repository"
)

// mockService implementa service.TaskService en memoria.
// Permite testear el handler sin tocar la capa de negocio ni la DB.
type mockService struct {
	tasks     []domain.Task
	task      domain.Task
	returnErr error
}

func (m *mockService) GetAll() ([]domain.Task, error)                               { return m.tasks, m.returnErr }
func (m *mockService) GetByID(id uint) (domain.Task, error)                         { return m.task, m.returnErr }
func (m *mockService) Create(req domain.CreateTaskRequest) (domain.Task, error)     { return m.task, m.returnErr }
func (m *mockService) Update(id uint, req domain.UpdateTaskRequest) (domain.Task, error) { return m.task, m.returnErr }
func (m *mockService) Delete(id uint) error                                          { return m.returnErr }

func setupRouter(svc *mockService) *chi.Mux {
	r := chi.NewRouter()
	h := handler.NewTaskHandler(svc)
	r.Route("/tasks", h.RegisterRoutes)
	return r
}

func TestGetAll_Exitoso(t *testing.T) {
	svc := &mockService{
		tasks: []domain.Task{
			{ID: 1, Title: "Tarea 1"},
			{ID: 2, Title: "Tarea 2"},
		},
	}
	r := setupRouter(svc)

	req := httptest.NewRequest("GET", "/tasks/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperaba status 200, obtuvo %d", rec.Code)
	}

	var tasks []domain.Task
	json.NewDecoder(rec.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Errorf("esperaba 2 tareas, obtuvo %d", len(tasks))
	}
}

func TestGetAll_ErrorServicio(t *testing.T) {
	svc := &mockService{returnErr: repository.ErrNotFound}
	r := setupRouter(svc)

	req := httptest.NewRequest("GET", "/tasks/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("esperaba status 500, obtuvo %d", rec.Code)
	}
}

func TestCreate_Exitoso(t *testing.T) {
	svc := &mockService{task: domain.Task{ID: 1, Title: "Nueva tarea"}}
	r := setupRouter(svc)

	body, _ := json.Marshal(domain.CreateTaskRequest{Title: "Nueva tarea"})
	req := httptest.NewRequest("POST", "/tasks/", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("esperaba status 201, obtuvo %d", rec.Code)
	}

	var task domain.Task
	json.NewDecoder(rec.Body).Decode(&task)
	if task.ID != 1 {
		t.Errorf("ID esperado 1, obtenido %d", task.ID)
	}
}

func TestCreate_CuerpoInvalido(t *testing.T) {
	svc := &mockService{}
	r := setupRouter(svc)

	req := httptest.NewRequest("POST", "/tasks/", bytes.NewReader([]byte("no es json")))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperaba status 400, obtuvo %d", rec.Code)
	}
}

func TestGetByID_Exitoso(t *testing.T) {
	svc := &mockService{task: domain.Task{ID: 1, Title: "Encontrada"}}
	r := setupRouter(svc)

	req := httptest.NewRequest("GET", "/tasks/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperaba status 200, obtuvo %d", rec.Code)
	}
}

func TestGetByID_NoEncontrado(t *testing.T) {
	svc := &mockService{returnErr: repository.ErrNotFound}
	r := setupRouter(svc)

	req := httptest.NewRequest("GET", "/tasks/999", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("esperaba status 404, obtuvo %d", rec.Code)
	}
}

func TestGetByID_IDInvalido(t *testing.T) {
	svc := &mockService{}
	r := setupRouter(svc)

	req := httptest.NewRequest("GET", "/tasks/abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperaba status 400, obtuvo %d", rec.Code)
	}
}

func TestUpdate_Exitoso(t *testing.T) {
	svc := &mockService{task: domain.Task{ID: 1, Title: "Actualizada", Completed: true}}
	r := setupRouter(svc)

	body, _ := json.Marshal(domain.UpdateTaskRequest{Title: "Actualizada", Completed: true})
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("esperaba status 200, obtuvo %d", rec.Code)
	}
}

func TestUpdate_NoEncontrado(t *testing.T) {
	svc := &mockService{returnErr: repository.ErrNotFound}
	r := setupRouter(svc)

	body, _ := json.Marshal(domain.UpdateTaskRequest{Title: "Fantasma"})
	req := httptest.NewRequest("PUT", "/tasks/999", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("esperaba status 404, obtuvo %d", rec.Code)
	}
}

func TestDelete_Exitoso(t *testing.T) {
	svc := &mockService{}
	r := setupRouter(svc)

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("esperaba status 204, obtuvo %d", rec.Code)
	}
}

func TestDelete_NoEncontrado(t *testing.T) {
	svc := &mockService{returnErr: repository.ErrNotFound}
	r := setupRouter(svc)

	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("esperaba status 404, obtuvo %d", rec.Code)
	}
}
