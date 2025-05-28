package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/Bharatjawa2/CtrlB_Assignment/internal/Storage"
	"github/Bharatjawa2/CtrlB_Assignment/internal/config"
	"github/Bharatjawa2/CtrlB_Assignment/models"
	"github/Bharatjawa2/CtrlB_Assignment/utils/response"
	"github/Bharatjawa2/CtrlB_Assignment/utils/security"
	"github/Bharatjawa2/CtrlB_Assignment/internal/http/middlewares"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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

func LoginStudent(storage storage.Storage, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		student, err := storage.GetStudentByEmail(creds.Email)
		if err != nil || !security.CheckPasswordHash(creds.Password, student.Password) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Generate JWT
		claims := jwt.MapClaims{
			"email": student.Email,
			"id":    student.Id,
			"role":  "student",  
			"exp":   time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		// Set token in cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    signedToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // true in production with HTTPS
			SameSite: http.SameSiteLaxMode,
		})

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "Login successful"})
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
		studentID, ok := r.Context().Value(middlewares.StudentIDKey).(int64)
		if !ok {
			http.Error(w, "Unauthorized: No student ID found", http.StatusUnauthorized)
			return
		}

		// Fetch existing student first
		existingStudent, err := storage.GetStudentById(studentID)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, map[string]string{"error": "Student not found"})
			return
		}

		var incomingData map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}

		// Update only fields present
		fullName, ok := incomingData["full_name"].(string);
		if ok {
			existingStudent.FullName = fullName
		}
		email, ok := incomingData["email"].(string);
		if ok {
			existingStudent.Email = email
		}
		password, ok := incomingData["password"].(string);
		if ok {
			existingStudent.Password = password
		}
		ageFloat, ok := incomingData["age"].(float64);
		if ok {
			existingStudent.Age = int(ageFloat)
		}
		gender, ok := incomingData["gender"].(string); 
		if ok {
			existingStudent.Gender = gender
		}
		phone, ok := incomingData["phone_number"].(string);
		if ok {
			existingStudent.PhoneNumber = phone
		}
		dob, ok := incomingData["dob"].(string);
		if ok {
			existingStudent.DOB = dob
		}
		address, ok := incomingData["address"].(string);
		if ok {
			existingStudent.Address = address
		}

		err = storage.UpdateStudent(studentID, existingStudent)
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

func Logout() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        http.SetCookie(w, &http.Cookie{
            Name:     "auth_token",
            Value:    "",
            Path:     "/",
            Expires:  time.Unix(0, 0), // Set expiration to Unix epoch to expire immediately
            HttpOnly: true,
            Secure:   false, // Set true in production with HTTPS
            SameSite: http.SameSiteLaxMode,
        })
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message":"Logged out successfully"}`))
    }
}
