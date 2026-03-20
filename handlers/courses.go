package handlers

import (
	"net/http"
	"strconv"

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

func ListCourses(c *gin.Context) {
	page, limit, offset := getPagination(c)

	var courses []models.Course
	var total int64

	config.DB.Model(&models.Course{}).Count(&total)
	config.DB.Limit(limit).Offset(offset).Find(&courses)

	c.JSON(http.StatusOK, gin.H{
		"data":  courses,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func getPagination(c *gin.Context) (page, limit, offset int) {
	page = 1
	limit = 20

	if p, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	offset = (page - 1) * limit
	return
}