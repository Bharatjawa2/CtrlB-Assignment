package storage

import "github/Bharatjawa2/CtrlB_Assignment/models"

type Storage interface{
	// Student
	CreateStudent(name string,email string,password string,age int,gender string,phoneNumber string,DOB string,Address string)(int64,error)
	GetStudentById(id int64) (models.Student,error)
	GetAllStudents() ([]models.Student, error)


	// Course
	CreateCourse(name string,description string,duration string,credits int,price int) (int64,error)

	// Enrollment
	EnrollStudent(studentID int64, courseID int64)(int64,error)
	GetCoursesByStudentID(studentID int64) ([]models.Course, error)
	GetStudentsByCourseID(courseID int64) ([]models.Student, error)
}