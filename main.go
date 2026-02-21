package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var db *gorm.DB
var jwtKey = []byte("secret_key")

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Role     string
}

type Course struct {
	ID    uint   `gorm:"primaryKey"`
	Title string
}

type Enrollment struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint
	CourseID uint
}

type Grade struct {
	ID           uint `gorm:"primaryKey"`
	EnrollmentID uint
	Score        float64
}

func main() {

	database, err := gorm.Open(sqlite.Open("grades.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = database

	// Migrate database FIRST
	db.AutoMigrate(&User{}, &Course{}, &Enrollment{}, &Grade{})

	// Seed users FIRST
	seedData()

	r := gin.Default()

	// Register routes BEFORE Run()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	r.POST("/login", login)

	// Run server LAST
	r.Run(":8080")
}

func seedData() {
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []User{
		{Name: "Admin", Email: "admin@test.com", Password: "123", Role: "admin"},
		{Name: "Teacher", Email: "teacher@test.com", Password: "123", Role: "teacher"},
		{Name: "Student", Email: "student@test.com", Password: "123", Role: "student"},
	}

	for _, u := range users {
		db.Create(&u)
	}
}

func generateToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func login(c *gin.Context) {
	var input struct {
		Email    string
		Password string
	}

	c.BindJSON(&input)

	var user User
	if err := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := generateToken(user)
	c.JSON(http.StatusOK, gin.H{"token": token})
}