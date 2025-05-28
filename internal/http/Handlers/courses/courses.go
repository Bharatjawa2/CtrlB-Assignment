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
