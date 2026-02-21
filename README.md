📘 Student Grade Management System API
📌 Overview
This project is a RESTful API built using Go (Golang) and the Gin framework that simulates a university-grade management system similar to Canvas or Blackboard.
The system supports:

Role-based authentication (Admin, Teacher, Student)
Course creation
Student enrollment
Grade assignment
GPA calculation
JWT-based secure access
Password hashing using bcrypt
The API demonstrates secure backend architecture, middleware usage, and structured database relationships.

🛠 Tech Stack

Go (Golang)
Gin Framework
GORM (ORM)
SQLite (Pure Go Driver – glebarez/sqlite)
JWT Authentication
bcrypt Password Hashing

🏗 Architecture Overview

Gin handles routing.
GORM manages database operations.
SQLite is used for lightweight local storage.
JWT enables stateless authentication.
Middleware enforces role-based authorization.
bcrypt secures user passwords.

🔐 Role Permissions Matrix
Role	Login	Create Course	Enroll Student	Assign Grade	View GPA
Admin	✅	      ✅	               ✅	        ❌	        ✅
Teacher	✅	      ❌	               ❌	        ✅	        ✅
Student	✅	      ❌	               ❌	        ❌	        ✅

🔁 API Endpoints
🔓 Public Routes
Health Check
GET /health
Login
POST /login

Body:

{
  "email": "admin@test.com",
  "password": "123"
}

🔒 Protected Routes (Require JWT Token)
All protected routes require header:
Authorization: Bearer <JWT_TOKEN>

Create Course (Admin Only)
POST /courses

Body:

{
  "title": "Operating Systems"
}

Enroll Student (Admin Only)
POST /enroll

Body:

{
  "user_id": 3,
  "course_id": 1
}

Assign Grade (Teacher Only)
POST /grades

Body:

{
  "enrollment_id": 1,
  "score": 88
}

View GPA (Admin / Teacher / Student)
GET /students/:id/performance

Example:

GET /students/3/performance

Response:

{
  "average_score": 88,
  "gpa": "3.52"
}

🧠 GPA Calculation Logic

Fetch all enrollments for the student.
Fetch grades linked to those enrollments.
Compute average score.
Convert to 4-point scale:
GPA = (average_score / 100) * 4

🔐 Security Features

JWT-based authentication
Role-based authorization using middleware
Password hashing using bcrypt
Input validation
Protected endpoints
Token expiration handling

🚀 Setup Instructions
1️⃣ Install Dependencies
    go mod tidy
2️⃣ Run the Application
    go run main.go
3️⃣ Start Fresh (Optional)

Delete the file:
grades.db
before running again for a clean database.

🧪 Complete Test Flow (Recommended Order)

Login as Admin
Create Course
Enroll Student
Login as Teacher
Assign Grade
View GPA

📂 Database Schema

User (Admin / Teacher / Student)
Course
Enrollment (User ↔ Course relationship)
Grade (linked to Enrollment)

🏁 Design Decisions

JWT chosen for stateless authentication.
Middleware ensures separation of authentication and business logic.
SQLite used for lightweight local execution.
bcrypt implemented to securely hash passwords.
Role-based access control mirrors real university hierarchy.

🔮 Future Improvements

Refresh token support
Pagination for large datasets
Unique constraint for preventing duplicate enrollments
PostgreSQL integration for production deployment
Logging & monitoring integration
Docker containerization
