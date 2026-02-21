# Design Document – Student Grade Management System API

## 1. System Overview

The system is a RESTful API designed to simulate a university grade management portal. It supports role-based access control for Admin, Teacher, and Student users.

## 2. Architecture

- Router: Gin Framework
- ORM: GORM
- Database: SQLite (pure Go driver)
- Authentication: JWT (stateless)
- Password Security: bcrypt hashing
- Middleware: Authentication and role-based authorization

## 3. Database Design

### Tables

User
- ID (Primary Key)
- Name
- Email (Unique)
- Password (Hashed)
- Role

Course
- ID
- Title

Enrollment
- ID
- UserID (Foreign Key)
- CourseID (Foreign Key)

Grade
- ID
- EnrollmentID (Foreign Key)
- Score

Relationships:
- One User can enroll in multiple Courses.
- One Course can have multiple Students.
- One Enrollment can have one Grade.

## 4. Authentication Flow

1. User logs in with email & password.
2. Password is verified using bcrypt.
3. JWT token is generated with:
   - user_id
   - role
   - expiration
4. Token must be included in Authorization header.

## 5. Authorization Strategy

Middleware extracts JWT claims and attaches:
- user_id
- role

Role middleware restricts endpoints based on allowed roles.

## 6. GPA Calculation

1. Fetch all enrollments for student.
2. Fetch grades for those enrollments.
3. Compute average score.
4. Convert to 4-point GPA scale.

GPA = (average_score / 100) * 4

## 7. Security Measures

- Password hashing
- Token expiration
- Role-based route protection
- Input validation