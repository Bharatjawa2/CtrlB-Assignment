package storage

import "github/Bharatjawa2/CtrlB_Assignment/models"

type Storage interface{
	CreateStudent(name string,email string,password string,age int,gender string,phoneNumber string,DOB string,Address string)(int64,error)
	GetStudentById(id int64) (models.Student,error)
}