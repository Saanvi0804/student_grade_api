package handlers

import (
	"fmt"
	"net/http"

	"grade-api/config"
	"grade-api/models"

	"github.com/gin-gonic/gin"
)

func EnrollStudent(c *gin.Context) {
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

	// Prevent duplicate enrollment
	var existing models.Enrollment
	if err := config.DB.Where("user_id = ? AND course_id = ?", input.UserID, input.CourseID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Student is already enrolled in this course"})
		return
	}

	enrollment := models.Enrollment{
		UserID:   input.UserID,
		CourseID: input.CourseID,
	}

	if result := config.DB.Create(&enrollment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll student"})
		return
	}

	c.JSON(http.StatusCreated, enrollment)
}

func GetPerformance(c *gin.Context) {
	studentID := c.Param("id")

	role := c.GetString("role")
	callerID := c.MustGet("user_id").(uint)
	if role == "student" && fmt.Sprintf("%d", callerID) != studentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Students can only view their own performance"})
		return
	}

	// Single JOIN query instead of N+1 DB calls
	type GradeRow struct {
		CourseTitle string
		Score       float64
	}

	var rows []GradeRow
	result := config.DB.
		Table("grades").
		Select("courses.title as course_title, grades.score").
		Joins("JOIN enrollments ON enrollments.id = grades.enrollment_id").
		Joins("JOIN courses ON courses.id = enrollments.course_id").
		Where("enrollments.user_id = ?", studentID).
		Scan(&rows)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve grades"})
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No grades found for this student"})
		return
	}

	var total float64
	type CourseGrade struct {
		Course string  `json:"course"`
		Score  float64 `json:"score"`
	}
	var breakdown []CourseGrade

	for _, row := range rows {
		total += row.Score
		breakdown = append(breakdown, CourseGrade{Course: row.CourseTitle, Score: row.Score})
	}

	avg := total / float64(len(rows))
	gpa := (avg / 100) * 4

	c.JSON(http.StatusOK, gin.H{
		"average_score": avg,
		"gpa":           fmt.Sprintf("%.2f", gpa),
		"breakdown":     breakdown,
	})
}