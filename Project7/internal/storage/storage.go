package storage

import "github.com/sahilsnghai/golang/Project7/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(age int) (types.Student, error)
	GetLists() ([]types.Student, error)
}
