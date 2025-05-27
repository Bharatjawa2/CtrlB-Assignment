package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/Bharatjawa2/CtrlB_Assignment/internal/Storage"
	"github/Bharatjawa2/CtrlB_Assignment/models"
	"github/Bharatjawa2/CtrlB_Assignment/utils/response"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a Student")
		var student models.Student
		err:=json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err,io.EOF) { // Body is Empty so ut give EOF Error
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}

		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}

		// request validation
		verr:=validator.New().Struct(student)
		if verr!=nil{
			validateError:=verr.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError(validateError))
			return
		}

		lastid,err:=storage.CreateStudent( 
			student.Name,
			student.Email,
			student.AGE,
		)

		slog.Info("User created successfully", slog.String("User Id: ",fmt.Sprint(lastid)))

		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,err)
		}



		response.WriteJson(w,http.StatusCreated,map[string]int64 {"id":lastid})
	}
}