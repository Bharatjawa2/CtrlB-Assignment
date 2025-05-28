package storage

import "github/Bharatjawa2/CtrlB_Assignment/models"

type Storage interface{
	// Student
	CreateStudent(name string,email string,password string,age int,gender string,phoneNumber string,DOB string,Address string)(int64,error)
	GetStudentByEmail(email string) (models.Student, error)
	GetStudentById(id int64) (models.Student,error)
	GetAllStudents() ([]models.Student, error)
	UpdateStudent(id int64, student models.Student) (error)


	// Course
	CreateCourse(name string,description string,duration string,credits int,price int) (int64,error)
	GetCourseById(id int64) (models.Course,error)
	GetAllCourses() ([]models.Course, error)
	UpdateCourse(id int64, course models.Course) (error)
	SearchCoursesByName(name string) ([]models.Course, error)

	// Enrollment
	EnrollStudent(studentID int64, courseID int64)(int64,error)
	GetCoursesByStudentID(studentID int64) ([]models.Course, error)
	GetStudentsByCourseID(courseID int64) ([]models.Student, error)
}