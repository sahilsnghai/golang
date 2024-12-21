package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sahilsnghai/golang/Project7/internal/types"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)

}

func GenrealError(err error) types.Response {

	return types.Response{Status: StatusError, Error: err.Error()}

}

func ValidationError(errs validator.ValidationErrors) types.Response {
	var errMsg []string

	for _, err := range errs {

		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is invalide", err.Field()))
		}

	}

	return types.Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ","),
	}

}
