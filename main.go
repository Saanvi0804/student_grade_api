package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

// ================== MAIN ==================

func main() {

	database, err := gorm.Open(sqlite.Open("grades.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = database

	db.AutoMigrate(&User{}, &Course{}, &Enrollment{}, &Grade{})
	seedData()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	r.POST("/login", login)

	protected := r.Group("/")
	protected.Use(authMiddleware())

	protected.GET("/protected", roleMiddleware("admin", "teacher", "student"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "You accessed a protected route"})
	})

	protected.POST("/courses", roleMiddleware("admin"), createCourse)
	protected.POST("/enroll", roleMiddleware("admin"), enrollStudent)
	protected.POST("/grades", roleMiddleware("teacher"), assignGrade)
	protected.GET("/students/:id/performance",
		roleMiddleware("admin", "teacher", "student"),
		getPerformance,
	)

	r.Run(":8080")
}

// ================== SEED DATA ==================

func seedData() {
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return
	}

	// Create users
	admin := User{Name: "Admin", Email: "admin@test.com", Password: hashPassword("123"), Role: "admin"}
	teacher := User{Name: "Teacher", Email: "teacher@test.com", Password: hashPassword("123"), Role: "teacher"}
	student := User{Name: "Student", Email: "student@test.com", Password: hashPassword("123"), Role: "student"}

	db.Create(&admin)
	db.Create(&teacher)
	db.Create(&student)

	// Create course
	course := Course{Title: "Mathematics"}
	db.Create(&course)

	// Enroll student properly using actual ID
	enrollment := Enrollment{
		UserID:   student.ID,
		CourseID: course.ID,
	}
	db.Create(&enrollment)

	// Assign grade
	grade := Grade{
		EnrollmentID: enrollment.ID,
		Score:        85,
	}
	db.Create(&grade)
}

// ================== AUTH ==================

func login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))

		c.Next()
	}
}

func roleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userRole := c.GetString("role")

		for _, role := range roles {
			if role == userRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
	}
}


func createCourse(c *gin.Context) {
	var course Course

	if err := c.BindJSON(&course); err != nil || course.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course title required"})
		return
	}

	db.Create(&course)
	c.JSON(http.StatusOK, course)
}

func enrollStudent(c *gin.Context) {
	var input struct {
		UserID   uint `json:"user_id"`
		CourseID uint `json:"course_id"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.UserID == 0 || input.CourseID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user or course ID"})
		return
	}

	enrollment := Enrollment{
		UserID:   input.UserID,
		CourseID: input.CourseID,
	}

	db.Create(&enrollment)
	c.JSON(http.StatusOK, enrollment)
}

func assignGrade(c *gin.Context) {
	var input struct {
		EnrollmentID uint    `json:"enrollment_id"`
		Score        float64 `json:"score"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.EnrollmentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid enrollment ID"})
		return
	}

	if input.Score < 0 || input.Score > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 0 and 100"})
		return
	}

	grade := Grade{
		EnrollmentID: input.EnrollmentID,
		Score:        input.Score,
	}

	db.Create(&grade)
	c.JSON(http.StatusOK, grade)
}


func getPerformance(c *gin.Context) {

	studentID := c.Param("id")

	var enrollments []Enrollment
	db.Where("user_id = ?", studentID).Find(&enrollments)

	var total float64
	var count int64

	for _, e := range enrollments {
		var grade Grade
		if err := db.Where("enrollment_id = ?", e.ID).First(&grade).Error; err == nil {
			total += grade.Score
			count++
		}
	}

	if count == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No grades found"})
		return
	}

	avg := total / float64(count)
	gpa := (avg / 100) * 4

	c.JSON(http.StatusOK, gin.H{
		"average_score": avg,
		"gpa":           fmt.Sprintf("%.2f", gpa),
	})
}


func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}