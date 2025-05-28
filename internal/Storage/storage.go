package storage

import "github/Bharatjawa2/CtrlB_Assignment/models"

type Storage interface{
	CreateStudent(name string,email string, age int)(int64,error)
	GetStudentById(id int64) (models.Student,error)
}