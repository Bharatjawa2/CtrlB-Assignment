package sqlite

import (
	"database/sql"
	"fmt"
	"github/Bharatjawa2/CtrlB_Assignment/internal/config"
	"github/Bharatjawa2/CtrlB_Assignment/models"
	security "github/Bharatjawa2/CtrlB_Assignment/utils/security"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

// Making Table

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// Create students table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	FullName TEXT,
	Email TEXT,
	Password TEXT,
	Age INTEGER,
	Gender TEXT,
	PhoneNumber TEXT,
	DOB TEXT,
	Address TEXT
	)`)

	if err != nil {
		return nil, fmt.Errorf("failed to create students table: %w", err)
	}

	// Create courses table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS courses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL,
		Description TEXT,
		Duration TEXT,
		Credits INTEGER,
		Price INTEGER
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create courses table: %w", err)
	}

	// Create Enrollment table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS enrollments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		student_id INTEGER,
		course_id INTEGER,
		enrolled_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(student_id) REFERENCES students(id),
		FOREIGN KEY(course_id) REFERENCES courses(id)
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create enrollments table: %w", err)
	}

	return &Sqlite{
		Db: db,
	}, nil
}

// Student

func (s *Sqlite) CreateStudent(FullName string, Email string, Password string, Age int, Gender string, PhoneNumber string, DOB string, Address string) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (FullName,Email,Password,Age,Gender,PhoneNumber,DOB,Address) VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(FullName, Email, Password, Age, Gender, PhoneNumber, DOB, Address)
	if err != nil {
		return 0, err
	}

	lastid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastid, nil
}

func (s *Sqlite) LoginStudent(email string, password string) (models.Student, error) {
    student, err := s.GetStudentByEmail(email)
    if err != nil {
        return models.Student{}, fmt.Errorf("invalid email or password")
    }
    if !security.CheckPasswordHash(password, student.Password) {
        return models.Student{}, fmt.Errorf("invalid email or password")
    }
    return student, nil
}


func (s *Sqlite) GetStudentByEmail(email string) (models.Student, error) {
    stmt, err := s.Db.Prepare("SELECT * FROM students WHERE Email = ? LIMIT 1")
    if err != nil {
        return models.Student{}, err
    }
    defer stmt.Close()

    var student models.Student
    err = stmt.QueryRow(email).Scan(
        &student.Id,
        &student.FullName,
        &student.Email,
        &student.Password,
        &student.Age,
        &student.Gender,
        &student.PhoneNumber,
        &student.DOB,
        &student.Address,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Student{}, fmt.Errorf("no Student Found with email %s", email)
        }
        return models.Student{}, err
    }

    return student, nil
}

func (s *Sqlite) GetStudentById(id int64) (models.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return models.Student{}, err
	}
	defer stmt.Close()

	var student models.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.FullName, &student.Email, &student.Password, &student.Age, &student.Gender, &student.PhoneNumber, &student.DOB, &student.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Student{}, fmt.Errorf("No Student Found with id %s", fmt.Sprint(id))
		}
		return models.Student{}, fmt.Errorf("Query Error: %w", err)
	}

	return student, nil
}

func (s *Sqlite) GetAllStudents() ([]models.Student, error) {
	rows, err := s.Db.Query("SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student
		err := rows.Scan(
			&student.Id,
			&student.FullName,
			&student.Email,
			&student.Password,
			&student.Age,
			&student.Gender,
			&student.PhoneNumber,
			&student.DOB,
			&student.Address,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (s *Sqlite) UpdateStudent(id int64, student models.Student) (error) {
	stmt, err := s.Db.Prepare(`UPDATE students SET 
		FullName = ?, 
		Email = ?, 
		Password = ?, 
		Age = ?, 
		Gender = ?, 
		PhoneNumber = ?, 
		DOB = ?, 
		Address = ? 
		WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		student.FullName,
		student.Email,
		student.Password, // We'll hash this later
		student.Age,
		student.Gender,
		student.PhoneNumber,
		student.DOB,
		student.Address,
		id,
	)
	return err
}

func (s *Sqlite) Logout()(error){
	return nil
}


// Courses

func (s *Sqlite) CreateCourse(Name string, Description string, Duration string, Credits int, Price int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO courses (Name,Description,Duration,Credits,Price) VALUES (?,?,?,?,?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(Name, Description, Duration, Credits, Price)
	if err != nil {
		return 0, err
	}

	lastid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastid, nil
}

func (s *Sqlite) GetCourseById(id int64) (models.Course, error) {
	stmt, err := s.Db.Prepare(`SELECT id, Name, Description, Duration, Credits, Price FROM courses WHERE id = ? LIMIT 1`)
	if err != nil {
		return models.Course{}, err
	}
	defer stmt.Close()

	var course models.Course
	err = stmt.QueryRow(id).Scan(&course.ID, &course.Name, &course.Description, &course.Duration, &course.Credits, &course.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Course{}, fmt.Errorf("no course found with id %d", id)
		}
		return models.Course{}, err
	}

	return course, nil
}

func (s *Sqlite) GetAllCourses() ([]models.Course, error) {
	rows, err := s.Db.Query("SELECT id, Name, Description, Duration, Credits, Price FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.Duration, &course.Credits, &course.Price)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func (s *Sqlite) UpdateCourse(id int64, course models.Course) error {
	stmt, err := s.Db.Prepare(`
		UPDATE courses SET 
			Name = ?, 
			Description = ?, 
			Duration = ?, 
			Credits = ?, 
			Price = ? 
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(course.Name, course.Description, course.Duration, course.Credits, course.Price, id)
	return err
}

func (s *Sqlite) SearchCoursesByName(name string) ([]models.Course, error) {
	query := `SELECT id, Name, Description, Duration, Credits, Price FROM courses WHERE Name LIKE ?`
	rows, err := s.Db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.Duration, &course.Credits, &course.Price)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	if len(courses) == 0 {
		return nil, fmt.Errorf("no courses found with name containing: %s", name)
	}

	return courses, nil
}

// Enrollment

func (s *Sqlite) EnrollStudent(studentID int64, courseID int64) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO enrollments (student_id, course_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(studentID, courseID)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (s *Sqlite) UnenrollStudent(studentID int64, courseID int64) error {
	stmt, err := s.Db.Prepare(`DELETE FROM enrollments WHERE student_id = ? AND course_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(studentID, courseID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no enrollment found for student %d in course %d", studentID, courseID)
	}

	return nil
}

func (s *Sqlite) GetStudentsByCourseID(courseID int64) ([]models.Student, error) {
	rows, err := s.Db.Query(`
		SELECT students.id, students.FullName, students.Email, students.Password, students.Age, 
		       students.Gender, students.PhoneNumber, students.DOB, students.Address
		FROM students
		INNER JOIN enrollments ON students.id = enrollments.student_id
		WHERE enrollments.course_id = ?
	`, courseID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.Id, &student.FullName, &student.Email, &student.Password, &student.Age,
			&student.Gender, &student.PhoneNumber, &student.DOB, &student.Address)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) GetCoursesByStudentID(studentID int64) ([]models.Course, error) {
	rows, err := s.Db.Query(`
		SELECT courses.id, courses.Name, courses.Description, courses.Duration, 
		       courses.Credits, courses.Price
		FROM courses
		INNER JOIN enrollments ON courses.id = enrollments.course_id
		WHERE enrollments.student_id = ?
	`, studentID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.Duration, &course.Credits, &course.Price)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}