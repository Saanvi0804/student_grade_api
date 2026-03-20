# Student Grade Management API

A RESTful API built with **Go**, **Gin**, **GORM**, and **SQLite** that simulates a university grade management system.

---

## Features

- JWT-based authentication with signing algorithm validation
- Role-based access control (Admin, Teacher, Student)
- bcrypt password hashing
- Pagination on list endpoints
- Students can only view their own performance
- Duplicate enrollment prevention
- Optimized database queries (single JOIN instead of N+1)
- Unit tested
- Dockerized with multi-stage build

---

## Tech Stack

| Layer      | Technology                        |
|------------|-----------------------------------|
| Language   | Go                                |
| Framework  | Gin                               |
| ORM        | GORM                              |
| Database   | SQLite                            |
| Auth       | JWT (golang-jwt/v5) + bcrypt      |
| Container  | Docker (multi-stage, alpine)      |

---

## Project Structure
```
.
├── main.go                  
├── config/
│   └── config.go            
├── handlers/
│   ├── auth.go              
│   ├── courses.go           
│   ├── students.go          
│   ├── grades.go            
│   └── handlers_test.go     
├── middleware/
│   └── middleware.go        
├── models/
│   └── models.go            
├── Dockerfile
├── .env.example
└── README.md
```

---

## Setup

### Environment Variables

Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

Generate a secure JWT secret:
```bash
openssl rand -hex 32
```

| Variable     | Required | Default          | Description                |
|-------------|----------|------------------|----------------------------|
| `JWT_SECRET` | Yes      | insecure default | Secret key for JWT signing |

### Run Locally
```bash
go mod tidy
go run main.go
```

### Run with Docker
```bash
docker build -t grade-api .
docker run -p 8080:8080 -e JWT_SECRET=your-secret grade-api
```

### Run Tests
```bash
go test ./handlers/...
```

---

## API Endpoints

### Public

| Method | Endpoint | Description   |
|--------|----------|---------------|
| GET    | /health  | Health check  |
| POST   | /login   | Get JWT token |

### Protected

All protected routes require: `Authorization: Bearer <JWT>`

| Method | Endpoint                    | Roles                   |
|--------|-----------------------------|-------------------------|
| GET    | /courses                    | admin, teacher, student |
| POST   | /courses                    | admin                   |
| GET    | /students                   | admin, teacher          |
| POST   | /enroll                     | admin                   |
| POST   | /grades                     | teacher                 |
| GET    | /students/:id/performance   | admin, teacher, student* |

*Students can only view their own performance.

List endpoints support pagination: `?page=1&limit=20`

---

## Role-Based Access Control

| Action            | Admin | Teacher | Student |
|------------------|-------|---------|---------|
| Create course     | ✅    | ❌      | ❌      |
| Enroll student    | ✅    | ❌      | ❌      |
| Assign grade      | ❌    | ✅      | ❌      |
| List students     | ✅    | ✅      | ❌      |
| List courses      | ✅    | ✅      | ✅      |
| View own GPA      | ✅    | ✅      | ✅      |
| View others' GPA  | ✅    | ✅      | ❌      |

---

## Seed Credentials

| Role    | Email              | Password |
|---------|--------------------|----------|
| Admin   | admin@test.com     | 123      |
| Teacher | teacher@test.com   | 123      |
| Student | student@test.com   | 123      |

---

## GPA Calculation

1. Fetch all enrollments for the student via a single JOIN query
2. Compute average score across all graded courses
3. Convert to 4.0 scale: `GPA = (average / 100) × 4`

Response includes a per-course breakdown:
```json
{
  "average_score": 85,
  "gpa": "3.40",
  "breakdown": [
    { "course": "Mathematics", "score": 85 }
  ]
}
```