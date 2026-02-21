# Student Grade Management System API

## Overview

This project is a RESTful API built using Go (Golang) and the Gin framework that simulates a university grade management system.

The system supports:

- Role-based authentication (Admin, Teacher, Student)
- Course creation
- Student enrollment
- Grade assignment
- GPA calculation
- JWT-based authentication
- Password hashing using bcrypt

This project demonstrates backend architecture design, middleware implementation, secure authentication, and relational database handling.

---

## Tech Stack

- Go (Golang)
- Gin Framework
- GORM
- SQLite (glebarez/sqlite driver)
- JWT (golang-jwt)
- bcrypt (password hashing)

---

## System Architecture

- Gin handles HTTP routing.
- GORM manages database interactions.
- SQLite provides lightweight local storage.
- JWT enables stateless authentication.
- Middleware enforces authentication and role-based authorization.
- bcrypt securely hashes passwords before storing them.

---

## Role-Based Access Control

| Role    | Login | Create Course | Enroll Student | Assign Grade | View GPA |
|---------|--------|---------------|----------------|--------------|----------|
| Admin   | Yes    | Yes           | Yes            | No           | Yes      |
| Teacher | Yes    | No            | No             | Yes          | Yes      |
| Student | Yes    | No            | No             | No           | Yes      |

---

## API Endpoints

### Public Routes

#### Health Check
GET /health

#### Login
POST /login

Example Request Body:

```json
{
  "email": "admin@test.com",
  "password": "123"
}
```

---

### Protected Routes

All protected routes require the following header:

Authorization: Bearer <JWT_TOKEN>

---

#### Create Course (Admin Only)
POST /courses

Example Body:

```json
{
  "title": "Operating Systems"
}
```

---

#### Enroll Student (Admin Only)
POST /enroll

Example Body:

```json
{
  "user_id": 3,
  "course_id": 1
}
```

---

#### Assign Grade (Teacher Only)
POST /grades

Example Body:

```json
{
  "enrollment_id": 1,
  "score": 88
}
```

---

#### View Student Performance
GET /students/:id/performance

Example:
GET /students/3/performance

Example Response:

```json
{
  "average_score": 88,
  "gpa": "3.52"
}
```

---

## GPA Calculation Logic

1. Fetch all enrollments for the student.
2. Fetch grades linked to those enrollments.
3. Compute the average score.
4. Convert percentage to a 4-point GPA scale.

Formula:

GPA = (average_score / 100) * 4

---

## Database Schema

### User
- ID (Primary Key)
- Name
- Email (Unique)
- Password (Hashed)
- Role (admin / teacher / student)

### Course
- ID
- Title

### Enrollment
- ID
- UserID (Foreign Key)
- CourseID (Foreign Key)

### Grade
- ID
- EnrollmentID (Foreign Key)
- Score

Relationships:
- One User can enroll in multiple Courses.
- One Course can have multiple Students.
- One Enrollment has one Grade.

---

## Authentication Flow

1. User logs in with email and password.
2. Password is verified using bcrypt.
3. A JWT token is generated containing:
   - user_id
   - role
   - expiration time
4. Token must be included in the Authorization header for protected routes.

---

## Security Features

- JWT-based stateless authentication
- Role-based authorization middleware
- Password hashing using bcrypt
- Input validation
- Token expiration handling
- Protected endpoints

---

## Setup Instructions

### Install Dependencies
go mod tidy

### Run the Application
go run main.go

### Reset Database (Optional)
Delete grades.db and restart the server.

---

## Recommended Test Flow

1. Login as Admin
2. Create Course
3. Enroll Student
4. Login as Teacher
5. Assign Grade
6. View GPA

---

## Design Decisions

- JWT chosen for stateless and scalable authentication.
- Middleware ensures separation of authentication and business logic.
- bcrypt used to securely hash passwords.
- SQLite selected for lightweight and zero-configuration setup.
- Role-based control models real academic hierarchy.

---

## Future Improvements

- Add composite unique constraints to prevent duplicate enrollments
- Add refresh token support
- Implement pagination for large datasets
- Migrate to PostgreSQL for scalability
- Add structured logging and monitoring
- Containerize using Docker

---

## Repository Contents

- main.go – API implementation
- README.md – Project documentation
- DESIGN.md – Detailed design explanation
- AI_PROMPTS.md – AI prompt transparency
