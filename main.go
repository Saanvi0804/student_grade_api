package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {

	database, err := gorm.Open(sqlite.Open("grades.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = database

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	r.Run(":8080")
}