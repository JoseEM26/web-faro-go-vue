package seed

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"proyecto-go/internal/domain"
)

const defaultPassword = "123456789"

func Run(db *gorm.DB) {
	seedUsers(db)
	seedCategories(db)
}

func seedUsers(db *gorm.DB) {
	var total int64
	db.Model(&domain.User{}).Count(&total)
	if total > 0 {
		log.Printf("[seed] %d usuario(s) ya existen, omitiendo seed de usuarios", total)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[seed] error generando hash: %v", err)
		return
	}

	emails := []string{
		"admin@taskgo.com",
		"juan.perez@taskgo.com",
		"maria.garcia@taskgo.com",
		"carlos.lopez@taskgo.com",
		"ana.martinez@taskgo.com",
	}

	for _, email := range emails {
		u := domain.User{Email: email, PasswordHash: string(hash)}
		if err := db.Create(&u).Error; err != nil {
			log.Printf("[seed] error creando %s: %v", email, err)
		}
	}

	log.Printf("[seed] %d usuarios creados — contrasena: %s", len(emails), defaultPassword)
}

func seedCategories(db *gorm.DB) {
	var total int64
	db.Model(&domain.Category{}).Count(&total)
	if total > 0 {
		return
	}

	cats := []domain.Category{
		{Name: "Trabajo",     Color: "#3B82F6"},
		{Name: "Personal",    Color: "#10B981"},
		{Name: "Urgente",     Color: "#EF4444"},
		{Name: "Aprendizaje", Color: "#8B5CF6"},
		{Name: "Ideas",       Color: "#F59E0B"},
	}

	for _, c := range cats {
		db.Create(&c)
	}
	log.Printf("[seed] %d categorias creadas", len(cats))
}
