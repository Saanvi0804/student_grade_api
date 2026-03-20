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

	var enrollments []models.Enrollment
	config.DB.Where("user_id = ?", studentID).Find(&enrollments)

	var total float64
	var count int64

	for _, e := range enrollments {
		var grade models.Grade
		if err := config.DB.Where("enrollment_id = ?", e.ID).First(&grade).Error; err == nil {
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