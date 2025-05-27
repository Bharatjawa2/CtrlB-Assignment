package student

import (
	"encoding/json"
	"errors"
	"github/Bharatjawa2/CtrlB_Assignment/models"
	"github/Bharatjawa2/CtrlB_Assignment/utils/response"
	"io"
	"log/slog"
	"net/http"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student models.Student
		err:=json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err,io.EOF) { // Body is Empty so ut give EOF Error
			response.WriteJson(w,http.StatusBadRequest,err.Error())
			return
		}
		slog.Info("Creating a Student")
		response.WriteJson(w,http.StatusCreated,map[string]string {"success":"OK"})
	}
}