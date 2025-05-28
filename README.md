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
   git clone <repository-url>
   cd admission-portal
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
   go run cmd/main.go
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
POST /api/v1/students
Content-Type: application/json

{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "1234567890",
    "address": "123 Main St"
}
```

#### Get Student
```http
GET /api/v1/students/{id}
```

#### Update Student
```http
PUT /api/v1/students/{id}
Content-Type: application/json

{
    "name": "John Doe Updated",
    "email": "john.updated@example.com"
}
```

#### Delete Student
```http
DELETE /api/v1/students/{id}
```

### Course Endpoints

#### Create Course
```http
POST /api/v1/courses
Content-Type: application/json

{
    "name": "Computer Science",
    "description": "Bachelor of Computer Science",
    "capacity": 100,
    "requirements": "High School Diploma"
}
```

#### Get Course
```http
GET /api/v1/courses/{id}
```

#### Update Course
```http
PUT /api/v1/courses/{id}
Content-Type: application/json

{
    "name": "Computer Science Updated",
    "capacity": 120
}
```

#### Delete Course
```http
DELETE /api/v1/courses/{id}
```

### Admission Endpoints

#### Submit Application
```http
POST /api/v1/admissions
Content-Type: application/json

{
    "student_id": 1,
    "course_id": 1,
    "status": "pending",
    "documents": ["transcript.pdf", "recommendation.pdf"]
}
```

#### Get Application
```http
GET /api/v1/admissions/{id}
```

#### Update Application Status
```http
PUT /api/v1/admissions/{id}/status
Content-Type: application/json

{
    "status": "approved"
}
```

## Database Schema

### Students Table
```sql
CREATE TABLE students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone TEXT,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Courses Table
```sql
CREATE TABLE courses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    capacity INTEGER NOT NULL,
    requirements TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Admissions Table
```sql
CREATE TABLE admissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    status TEXT NOT NULL,
    documents TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id),
    FOREIGN KEY (course_id) REFERENCES courses(id)
);
```

## Testing

Run the test suite:
```bash
go test ./...
```

## Docker Support

The application includes Docker support for easy deployment:

1. Build the Docker image:
   ```bash
   docker build -t admission-portal .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 admission-portal
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
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Address   string    `json:"address"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Course Model
```go
type Course struct {
    ID           int64     `json:"id"`
    Name         string    `json:"name"`
    Description  string    `json:"description"`
    Capacity     int       `json:"capacity"`
    Requirements string    `json:"requirements"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### Admission Model
```go
type Admission struct {
    ID        int64     `json:"id"`
    StudentID int64     `json:"student_id"`
    CourseID  int64     `json:"course_id"`
    Status    string    `json:"status"`
    Documents []string  `json:"documents"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
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

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
