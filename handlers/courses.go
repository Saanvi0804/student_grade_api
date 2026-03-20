package handlers

import (
	"net/http"

	"grade-api/config"
	"grade-api/models"

	"github.com/gin-gonic/gin"
)

func CreateCourse(c *gin.Context) {
	var course models.Course

	if err := c.BindJSON(&course); err != nil || course.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course title required"})
		return
	}

	if result := config.DB.Create(&course); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}