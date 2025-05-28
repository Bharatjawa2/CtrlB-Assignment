# Admission Portal Backend

A robust backend system for managing student admissions, courses, and enrollment processes. Built with Go, this system provides a RESTful API for handling all admission-related operations.

## Table of Contents
- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Testing](#testing)
- [Docker Support](#docker-support)
- [Models](#models)

## Features

- **Student Management**
  - Create, read, update, and delete student profiles
  - Track student status and enrollment history
  - Manage student documents and information

- **Course Management**
  - Create and manage courses
  - Set course capacity and requirements
  - Track course enrollment

- **Admission Process**
  - Handle admission applications
  - Process admission decisions
  - Manage enrollment status
  - Track application status

## Project Structure

```
.
├── cmd/                    # Application entry points (main packages for various apps or services)
├── config/                 # Configuration files (e.g., YAML, JSON, ENV, or Go configs)
├── internal/               # Private application code (only importable within this module)
│   ├── config/             # Internal config-related logic (parsing, loading, validation)
│   ├── http/               # HTTP handlers and routers
│       ├── Handlers/
│           ├── admin       # HTTP handlers for admin-related endpoints
│           ├── courses     # HTTP handlers for courses-related endpoints
│           ├── enrollment  # HTTP handlers for enrollment-related endpoints
│           ├── student     # HTTP handlers for student-related endpoints
│   └── middleware/         # HTTP middleware (auth, logging, recovery, etc.)
├── models/                 # Data models (structs representing database entities or API payloads)
├── storage/                # Database and storage implementations (e.g., SQL queries, ORM, file storage)
├── utils/                  # Utility functions (helpers shared across the app)
├── Dockerfile              # Docker image build instructions
├── docker-compose.yml      # Docker Compose setup for multi-container environments
├── go.mod                  # Go module definition (dependencies and module path)
└── go.sum                  # Go module checksum file (verifies dependency integrity)

```

## Prerequisites

- Go 1.21 or higher
- SQLite3
- Docker (optional, for containerized deployment)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Bharatjawa2/CtrlB-Assignment
   cd CTRLB
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Configuration

The application uses a configuration file located at `config/config.yaml`. You can modify the following settings:

```yaml
server:
  port: 8082
  host: localhost

database:
  type: sqlite
  path: ./storage/storage.db

logging:
  level: info
  format: json
```

## Running the Application

### Local Development

1. Start the application:
   ```bash
   go run cmd/CTRLB/main.go
   ```

2. The server will start on `http://localhost:8082`

### Using Docker

1. Build and run using Docker Compose:
   ```bash
   docker-compose up --build
   ```

## API Documentation

### Student Endpoints

#### Create Student
```http
POST /api/students
Content-Type: application/json

{
    {
    "full_name": "Test2",
    "Email": "test2@gmail.com",
    "Password": "testing123",
    "Age": 22,
    "gender": "male",
    "Phone_Number": "9812345644",
    "DOB": "2003-05-37",
    "Address": "Indore"
  }
}
```

#### Get Student
```http
GET /api/students/{id}
GET /api/students/all
```

#### Update Student
```http
PUT /api/students/update

{
    {
    "ID": 10,
    "full_name": "Pooja",
    "Email": "pooja.verma@gmail.com",
    "Password": "Pooja@2022",
    "Age": 22,
    "gender": "female",
    "Phone_Number": "9090909090",
    "DOB": "2002-10-14",
    "Address": "Pune"
  }
}
```

### Course Endpoints

#### Create Course
```http
POST /api/courses

{
    {
    "ID": 10,
    "Name": "Business Communication",
    "Description": "Focuses on professional communication, presentation skills, business writing, and interpersonal effectiveness in corporate environments.",
    "Duration": "2 months",
    "Credits": 2,
    "Price": 690
  }
}
```

#### Get Course
```http
GET /api/courses/{id}
GET /api/courses/all
```

#### Update Course
```http
PUT /api/courses/update/{id}

{
     {
    "ID": 8,
    "Name": "Graphic Design Basics",
    "Description": "A practical course on design principles, color theory, typography, and digital tools like Adobe Photoshop and Illustrator.",
    "Duration": "4 months",
    "Credits": 3,
    "Price": 1099
  },
}
```

### Enrollment Endpoints

#### Enroll in Course
```http
POST /api/enrollment

{
    "student_id": 1,
    "course_id": 1,
}
```

#### Get Enrolled Students and Courses
```http
GET /api/enrolled/students/{id}
GET /api/enrolled/courses/{id}
```

#### Remove Enrollment
```http
POST /api/unenrollment
```

## Docker Support

The application includes Docker support for easy deployment:

1. Build the Docker image:
   ```bash
   docker-compose up --build
   ```

2. Run the container:
   ```bash
    docker run -p 8082:8082 ctrlb-backend
   ```

Or use Docker Compose:
```bash
docker-compose up
```

## Models

The application uses the following data models to represent the core entities:

### Student Model
```go
type Student struct {
	Id              int64 
	FullName        string 
	Email           string 
	Password        string 
	Age             int    
	Gender          string 
	PhoneNumber     string 
	DOB             string 
	Address         string 
}

```

### Course Model
```go
type Course struct {
    ID          int64  
    Name        string 
    Description string 
    Duration    string  
    Credits     int 
	Price		int	   
}

```

### Enrollment Model
```go
type Enrollment struct {
	ID        int64 
	StudentID int64 
	CourseID  int64 
}

```

Each model includes:
- Primary key (ID)
- Relevant fields for the entity
- Timestamps for creation and updates
- JSON tags for serialization
- Appropriate data types for each field

The models are designed to:
- Maintain data integrity
- Support CRUD operations
- Enable proper data validation
- Facilitate API responses
- Support database operations


