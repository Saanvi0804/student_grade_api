package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `gorm:"password" json:"-"`
	Role     string `json:"role"`
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