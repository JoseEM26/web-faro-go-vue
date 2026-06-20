//go:build integration

// Los tests de repositorio requieren una base de datos real.
// Ejecutar con: go test -tags=integration ./internal/repository/...
// La variable TEST_DATABASE_URL puede apuntar a una DB de test dedicada.
package repository_test

import (
	"errors"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"proyecto-go/internal/domain"
	"proyecto-go/internal/repository"
)

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=proyecto_go_test sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("conectando a la DB de test: %v", err)
	}

	if err := db.AutoMigrate(&domain.Task{}); err != nil {
		t.Fatalf("ejecutando migraciones: %v", err)
	}

	t.Cleanup(func() {
		db.Exec("DELETE FROM tasks")
	})

	return db
}

func TestGorm_Create(t *testing.T) {
	repo := repository.NewGormTaskRepository(setupDB(t))

	task := &domain.Task{Title: "Test GORM", Description: "desc"}
	if err := repo.Create(task); err != nil {
		t.Fatalf("error al crear: %v", err)
	}
	if task.ID == 0 {
		t.Error("se esperaba un ID generado por la DB")
	}
}

func TestGorm_FindAll(t *testing.T) {
	repo := repository.NewGormTaskRepository(setupDB(t))

	repo.Create(&domain.Task{Title: "Tarea A"})
	repo.Create(&domain.Task{Title: "Tarea B"})

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(tasks) < 2 {
		t.Errorf("esperaba al menos 2 tareas, obtuvo %d", len(tasks))
	}
}

func TestGorm_FindByID(t *testing.T) {
	repo := repository.NewGormTaskRepository(setupDB(t))

	created := &domain.Task{Title: "Buscar por ID"}
	repo.Create(created)

	found, err := repo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("error al buscar: %v", err)
	}
	if found.Title != "Buscar por ID" {
		t.Errorf("titulo esperado %q, obtenido %q", "Buscar por ID", found.Title)
	}
}

func TestGorm_FindByID_NoEncontrado(t *testing.T) {
	repo := repository.NewGormTaskRepository(setupDB(t))

	_, err := repo.FindByID(99999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("esperaba ErrNotFound, obtuvo %v", err)
	}
}

func TestGorm_Update(t *testing.T) {
	repo := repository.NewGormTaskRepository(setupDB(t))

	task := &domain.Task{Title: "Original"}
	repo.Create(task)
	task.Title = "Actualizada"
	task.Completed = true

	if err := repo.Update(task); err != nil {
		t.Fatalf("error al actualizar: %v", err)
	}

	found, _ := repo.FindByID(task.ID)
	if found.Title != "Actualizada" {
		t.Errorf("titulo esperado %q, obtenido %q", "Actualizada", found.Title)
	}
}

func TestGorm_Delete(t *testing.T) {
	repo := repository.NewGormTaskRepository(setupDB(t))

	task := &domain.Task{Title: "A eliminar"}
	repo.Create(task)

	if err := repo.Delete(task.ID); err != nil {
		t.Fatalf("error al eliminar: %v", err)
	}

	_, err := repo.FindByID(task.ID)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Error("se esperaba ErrNotFound tras eliminar")
	}
}
