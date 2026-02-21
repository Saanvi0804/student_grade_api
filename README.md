📘 Student Grade Management System API
📌 Project Overview

This project is a RESTful API built using Go (Golang) and the Gin framework to simulate a university-grade management system similar to Canvas or Blackboard.

The system supports:

Role-based authentication (Admin, Teacher, Student)

Course creation

Student enrollment

Grade assignment

GPA calculation

Secure JWT-based authentication

Password hashing using bcrypt

This project demonstrates secure backend architecture, middleware implementation, role-based access control, and relational database handling.

🛠 Tech Stack

Language: Go (Golang)

Framework: Gin

ORM: GORM

Database: SQLite (Pure Go Driver – glebarez/sqlite)

Authentication: JWT (golang-jwt)

Password Security: bcrypt

🏗 System Architecture

Gin handles HTTP routing.

GORM manages database interactions.

SQLite provides lightweight local storage.

JWT enables stateless authentication.

Middleware enforces authentication and role-based authorization.

bcrypt securely hashes user passwords before storage.

🔐 Role-Based Access Control
Role	Login	Create Course	Enroll Student	Assign Grade	View GPA
Admin	✅	✅	✅	❌	✅
Teacher	✅	❌	❌	✅	✅
Student	✅	❌	❌	❌	✅
🔁 API Endpoints
🔓 Public Routes
Health Check
GET /health
Login
POST /login

Example Body:

{
  "email": "admin@test.com",
  "password": "123"
}
🔒 Protected Routes

All protected routes require:

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
View Student GPA
GET /students/:id/performance

Example:

GET /students/3/performance

Example Response:

{
  "average_score": 88,
  "gpa": "3.52"
}
🧮 GPA Calculation Logic

Fetch all enrollments for the student.

Fetch grades linked to those enrollments.

Compute the average score.

Convert percentage to 4-point scale.

Formula:

GPA = (average_score / 100) * 4
🔐 Security Features

JWT-based stateless authentication

Role-based authorization middleware

Password hashing using bcrypt

Input validation

Token expiration handling

Protected routes

🗄 Database Schema
User

ID (Primary Key)

Name

Email (Unique)

Password (Hashed)

Role (admin / teacher / student)

Course

ID

Title

Enrollment

ID

UserID (Foreign Key)

CourseID (Foreign Key)

Grade

ID

EnrollmentID (Foreign Key)

Score

🚀 Setup Instructions
1️⃣ Install Dependencies
go mod tidy
2️⃣ Run the Application
go run main.go
3️⃣ Start Fresh (Optional)

Delete:

grades.db

before running again to reset the database.

🧪 Recommended Test Flow

Login as Admin

Create Course

Enroll Student

Login as Teacher

Assign Grade

View GPA

🧠 Design Decisions

JWT chosen for stateless and scalable authentication.

Middleware ensures separation of authentication and business logic.

bcrypt used to securely hash passwords.

SQLite selected for lightweight, zero-configuration setup.

Role-based control models real academic hierarchy.

🔮 Future Improvements

Prevent duplicate enrollments using composite unique constraints

Add refresh token support

Implement pagination for large datasets

Migrate to PostgreSQL for production scalability

Add logging and monitoring

Containerize using Docker

📂 Repository Contents

main.go — Complete API implementation

README.md — Project overview and usage

DESIGN.md — Detailed design documentation

AI_PROMPTS.md — Transparency of AI-assisted development
