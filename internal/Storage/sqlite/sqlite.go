package sqlite

import (
	"database/sql"
	"fmt"
	"github/Bharatjawa2/CtrlB_Assignment/internal/config"
	"github/Bharatjawa2/CtrlB_Assignment/models"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct{
	Db *sql.DB
}


// Students 

func New(cfg *config.Config)(*Sqlite,error){
	db,err:=sql.Open("sqlite3",cfg.StoragePath)
	if err!=nil{
		return nil,err
	}
	
	// Create students table
	_,err=db.Exec(`CREATE TABLE IF NOT EXISTS students (
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

	if err!=nil{
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

	return &Sqlite{
		Db:db,
	},nil
}

func (s *Sqlite) CreateStudent(FullName string,Email string,Password string, Age int,Gender string,PhoneNumber string,DOB string,Address string)(int64,error){
	stmt,err:=s.Db.Prepare("INSERT INTO students (FullName,Email,Password,Age,Gender,PhoneNumber,DOB,Address) VALUES (?,?,?,?,?,?,?,?)")
	if err!=nil{
		return 0,err
	}

	defer stmt.Close()

	result,err:=stmt.Exec(FullName,Email,Password,Age,Gender,PhoneNumber,DOB,Address)
	if err!=nil{
		return 0,err
	}

	lastid,err:=result.LastInsertId()
	if err!=nil{
		return 0,err
	}

	return lastid ,nil
}


func (s *Sqlite) GetStudentById(id int64) (models.Student,error){
	stmt,err:=s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err!=nil{
		return models.Student{},err
	}
	defer stmt.Close()

	var student models.Student

	err=stmt.QueryRow(id).Scan(&student.Id, &student.FullName, &student.Email, &student.Password, &student.Age, &student.Gender,& student.PhoneNumber,&student.DOB,&student.Address)
	if err!=nil{
		if err==sql.ErrNoRows{
			return models.Student{},fmt.Errorf("No Student Found with id %s",fmt.Sprint(id))
		}
		return models.Student{},fmt.Errorf("Query Error: %w",err)
	}

	return student,nil
}



// Courses

func (s *Sqlite) CreateCourse(Name string,Description string,Duration string, Credits int, Price int)(int64,error){
	stmt,err:=s.Db.Prepare("INSERT INTO courses (Name,Description,Duration,Credits,Price) VALUES (?,?,?,?,?)")
	if err!=nil{
		return 0,err
	}

	defer stmt.Close()

	result,err:=stmt.Exec(Name,Description,Duration,Credits,Price)
	if err!=nil{
		return 0,err
	}

	lastid,err:=result.LastInsertId()
	if err!=nil{
		return 0,err
	}

	return lastid ,nil
}