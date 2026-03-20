package config

import (
	"log"
	"os"

	"grade-api/models"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	database, err := gorm.Open(sqlite.Open("grades.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	DB = database
	DB.AutoMigrate(&models.User{}, &models.Course{}, &models.Enrollment{}, &models.Grade{})
	seedData()
}

func JWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET not set. Set this in production!")
		secret = "change-me-in-production"
	}
	return []byte(secret)
}

func seedData() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	admin := models.User{Name: "Admin", Email: "admin@test.com", Password: hashPassword("123"), Role: "admin"}
	teacher := models.User{Name: "Teacher", Email: "teacher@test.com", Password: hashPassword("123"), Role: "teacher"}
	student := models.User{Name: "Student", Email: "student@test.com", Password: hashPassword("123"), Role: "student"}

	DB.Create(&admin)
	DB.Create(&teacher)
	DB.Create(&student)

	course := models.Course{Title: "Mathematics"}
	DB.Create(&course)

	enrollment := models.Enrollment{UserID: student.ID, CourseID: course.ID}
	DB.Create(&enrollment)

	grade := models.Grade{EnrollmentID: enrollment.ID, Score: 85}
	DB.Create(&grade)

	log.Println("Seed data created.")
}

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("failed to hash password:", err)
	}
	return string(hashed)
}