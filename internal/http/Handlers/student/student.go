package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/Bharatjawa2/CtrlB_Assignment/internal/Storage"
	"github/Bharatjawa2/CtrlB_Assignment/models"
	"github/Bharatjawa2/CtrlB_Assignment/utils/response"
	"github/Bharatjawa2/CtrlB_Assignment/utils/security"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func Register(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a Student")
		var student models.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) { // Body is Empty so ut give EOF Error
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		verr := validator.New().Struct(student)
		if verr != nil {
			validateError := verr.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateError))
			return
		}

		hashedPassword, err:=security.HashPassword(student.Password)
		if err!=nil{
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		lastid, err := storage.CreateStudent(
			student.FullName,
			student.Email,
			hashedPassword,
			student.Age,
			student.Gender,
			student.PhoneNumber,
			student.DOB,
			student.Address,
		)

		slog.Info("User created successfully", slog.String("User Id: ", fmt.Sprint(lastid)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastid})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))
		Intid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(Intid)
		if err != nil {
			slog.Info("Error getting user", slog.String("Id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetAllStudents(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Fetching all students")

		students, err := storage.GetAllStudents()
		if err!=nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Missing student ID"})
			return
		}

		studentId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Fetch existing student first
		existingStudent, err := storage.GetStudentById(studentId)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, map[string]string{"error": "Student not found"})
			return
		}

		// Decode incoming JSON into a map to check which fields are provided
		var incomingData map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}

		// Update only fields that are present in incomingData
		if fullName, ok := incomingData["full_name"].(string); ok {
			existingStudent.FullName = fullName
		}
		if email, ok := incomingData["email"].(string); ok {
			existingStudent.Email = email
		}
		if password, ok := incomingData["password"].(string); ok {
			existingStudent.Password = password
		}
		if ageFloat, ok := incomingData["age"].(float64); ok {
			existingStudent.Age = int(ageFloat) // JSON numbers come as float64
		}
		if gender, ok := incomingData["gender"].(string); ok {
			existingStudent.Gender = gender
		}
		if phone, ok := incomingData["phone_number"].(string); ok {
			existingStudent.PhoneNumber = phone
		}
		if dob, ok := incomingData["dob"].(string); ok {
			existingStudent.DOB = dob
		}
		if address, ok := incomingData["address"].(string); ok {
			existingStudent.Address = address
		}

		// Now update student in storage
		err = storage.UpdateStudent(studentId, existingStudent)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update student"})
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "Student updated successfully"})
	}
}

func GetStudentByEmail(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email") // Want to Fetch from Query parameters 
		if email == "" {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Missing email query parameter"})
			return
		}
		student, err := storage.GetStudentByEmail(email)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, map[string]string{"error": "Student not found"})
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}
