package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sahilsnghai/golang/Project7/internal/storage"
	"github.com/sahilsnghai/golang/Project7/internal/types"
	"github.com/sahilsnghai/golang/Project7/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a new student")

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GenrealError(fmt.Errorf("empty body")))
			return

		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenrealError(err))
		}

		if err := validator.New().Struct(student); err != nil {
			validateError := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateError))
			return
		}

		lastID, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}
		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastID)))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastID})
	}
}

func GetbyId(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		_id, err := strconv.Atoi(id)
		if err != nil {
			slog.Error("converting string id into int ")
			response.WriteJson(w, http.StatusBadRequest, response.GenrealError(err))
			return
		}

		student, err := storage.GetStudentById(_id)

		if err != nil {
			slog.Error("error getting in user")
			response.WriteJson(w, http.StatusInternalServerError, response.GenrealError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)

	}
}

func GetLists(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		students, err := storage.GetLists()

		if err != nil {
			slog.Error("error getting students data", slog.String("error ", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GenrealError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, students)

	}
}
