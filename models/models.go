package models

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