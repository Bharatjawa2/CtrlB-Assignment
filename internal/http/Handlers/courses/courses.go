package courses

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

func CreateCourse(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Creating a Course")
		var course models.Course
		err:=json.NewDecoder(r.Body).Decode(&course)

		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validator
		verr:=validator.New().Struct(course)
		if verr!=nil{
			validatorError:=verr.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError((validatorError)))
			return
		}

		lastCourseId,err:=storage.CreateCourse(
			course.Name,
			course.Description,
			course.Duration,
			course.Credits,
			course.Price,
		)

		slog.Info("Courses Added successfully",slog.String("Course ID: ",fmt.Sprint(lastCourseId)))
		if err!=nil{
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastCourseId})
	}
}

func GetCourseById(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		id:=r.PathValue("id")
		slog.Info("Getting a Course",slog.String("id",id))
		CourseId,err:=strconv.ParseInt(id,10,64)
		if err!=nil{
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		course,err:=storage.GetCourseById(CourseId)
		if err!=nil{
			slog.Info("Error getting course", slog.String("Id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK,course)
	}
}

func GetAllCourses(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		courses, err := storage.GetAllCourses()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch courses"})
			return
		}
		response.WriteJson(w, http.StatusOK, courses)
	}
}

func UpdateCourse(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Missing course ID"})
			return
		}

		courseId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Fetch existing course first
		existingCourse, err := storage.GetCourseById(courseId)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, map[string]string{"error": "Course not found"})
			return
		}

		// Decode incoming JSON into a map
		var incomingData map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
			return
		}

		// Update only fields present in JSON
		if name, ok := incomingData["name"].(string); ok {
			existingCourse.Name = name
		}
		if desc, ok := incomingData["description"].(string); ok {
			existingCourse.Description = desc
		}
		if durationStr, ok := incomingData["duration"].(string); ok {
			existingCourse.Duration = durationStr
		}
		if creditsFloat, ok := incomingData["credits"].(float64); ok {
			existingCourse.Credits = int(creditsFloat)
		}
		if priceFloat, ok := incomingData["price"].(float64); ok {
			existingCourse.Price = int(priceFloat)
		}

		// Save updated course
		err = storage.UpdateCourse(courseId, existingCourse)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update course"})
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "Course updated successfully"})
	}
}

func SearchCoursesByName(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Missing course name query parameter"})
			return
		}

		courses, err := storage.SearchCoursesByName(name)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}

		response.WriteJson(w, http.StatusOK, courses)
	}
}
