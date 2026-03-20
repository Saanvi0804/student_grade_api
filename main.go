package main

import (
	"net/http"

	"grade-api/config"
	"grade-api/handlers"
	"grade-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/login", handlers.Login)

	protected := r.Group("/")
	protected.Use(middleware.Auth())

	protected.POST("/courses", middleware.RequireRole("admin"), handlers.CreateCourse)
	protected.POST("/enroll", middleware.RequireRole("admin"), handlers.EnrollStudent)
	protected.POST("/grades", middleware.RequireRole("teacher"), handlers.AssignGrade)
	protected.GET("/students/:id/performance",
		middleware.RequireRole("admin", "teacher", "student"),
		handlers.GetPerformance,
	)

	r.Run(":8080")
}