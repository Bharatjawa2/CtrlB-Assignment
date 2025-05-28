package enrollment

import (
	"encoding/json"
	"errors"
	"fmt"
	storage "github/Bharatjawa2/CtrlB_Assignment/internal/Storage"
	"github/Bharatjawa2/CtrlB_Assignment/models"
	"github/Bharatjawa2/CtrlB_Assignment/utils/response"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func EnrollStudent(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Enrolling Student")
		var enroll models.Enrollment
		err:=json.NewDecoder(r.Body).Decode(&enroll)

		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validator
		verr:=validator.New().Struct(enroll)
		if verr!=nil{
			validatorError:=verr.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError((validatorError)))
			return
		}

		lastEnrollId,err:=storage.EnrollStudent(
			enroll.StudentID,
			enroll.CourseID,
		)
		slog.Info("Student Enrolled Successfully",slog.String("Enrollment ID: ",fmt.Sprint(lastEnrollId)))
		if err!=nil{
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastEnrollId})
	
	}
}

func GetCoursesByStudentID(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		slog.Info("List of Courses Registered by Student")
		
		// Extract student ID from URL
		id := r.PathValue("id")
		studentId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		
		// Fetch courses by student ID
		courses, err := storage.GetCoursesByStudentID(studentId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, courses)

	}
}

func GetStudentsByCourseID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("List of Students Registered in Course")

		// Extract course ID from URL
		id := r.PathValue("id")
		courseID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Fetch students by course ID
		students, err := storage.GetStudentsByCourseID(courseID)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func UnenrollStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var enrollment models.Enrollment

		err := json.NewDecoder(r.Body).Decode(&enrollment)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}

		if enrollment.StudentID == 0 || enrollment.CourseID == 0 {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Missing student_id or course_id"})
			return
		}

		err = storage.UnenrollStudent(enrollment.StudentID, enrollment.CourseID)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "Student unenrolled from course successfully"})
	}
}
