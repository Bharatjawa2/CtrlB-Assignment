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