package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"proyecto-go/internal/config"
	"proyecto-go/internal/domain"
	"proyecto-go/internal/handler"
	"proyecto-go/internal/middleware"
	"proyecto-go/internal/repository"
	"proyecto-go/internal/seed"
	"proyecto-go/internal/service"
	"proyecto-go/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("cargando configuracion: %v", err)
	}

	db, err := database.NewPostgresDB(cfg.DSN())
	if err != nil {
		log.Fatalf("conectando a la base de datos: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Task{}, &domain.Category{}, &domain.Note{}); err != nil {
		log.Fatalf("ejecutando migraciones: %v", err)
	}

	seed.Run(db)

	// Repositorios
	userRepo     := repository.NewGormUserRepository(db)
	taskRepo     := repository.NewGormTaskRepository(db)
	categoryRepo := repository.NewGormCategoryRepository(db)
	noteRepo     := repository.NewGormNoteRepository(db)

	// Servicios
	authSvc     := service.NewAuthService(userRepo, cfg.JWTSecret)
	taskSvc     := service.NewTaskService(taskRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	noteSvc     := service.NewNoteService(noteRepo)

	// Handlers
	authH     := handler.NewAuthHandler(authSvc)
	taskH     := handler.NewTaskHandler(taskSvc)
	categoryH := handler.NewCategoryHandler(categorySvc)
	noteH     := handler.NewNoteHandler(noteSvc)

	r := chi.NewRouter()
	r.Use(middleware.CORS)
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", authH.RegisterRoutes)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWTSecret))
			r.Route("/tasks",      taskH.RegisterRoutes)
			r.Route("/categories", categoryH.RegisterRoutes)
			r.Route("/notes",      noteH.RegisterRoutes)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Servidor corriendo en http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
