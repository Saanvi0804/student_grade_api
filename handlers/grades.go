package handlers

import (
	"net/http"

	"grade-api/config"
	"grade-api/models"

	"github.com/gin-gonic/gin"
)

func AssignGrade(c *gin.Context) {
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

	grade := models.Grade{
		EnrollmentID: input.EnrollmentID,
		Score:        input.Score,
	}

	if result := config.DB.Create(&grade); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign grade"})
		return
	}

	c.JSON(http.StatusCreated, grade)
}