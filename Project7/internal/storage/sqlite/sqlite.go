package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sahilsnghai/golang/Project7/internal/config"
	"github.com/sahilsnghai/golang/Project7/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	fmt.Println(cfg.StoragePath)
	Db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER NOT NULL,
    email TEXT NOT NULL

	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: Db}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	slog.Info("Created Student called")
	stmt, err := s.Db.Prepare("Insert into students (name, email, age) values (?,?,?)")
	fmt.Printf("name email age: %s %s %d \n", name, email, age)

	if err != nil {
		slog.Error("found error while perparing %s", slog.String("error", err.Error()))
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		slog.Error("found error while exec ", slog.String("error", err.Error()))

		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		slog.Error("found error while fetching last id ", slog.String("error", err.Error()))

		return 0, err
	}

	return lastId, err
}

func (s *Sqlite) GetStudentById(id int) (types.Student, error) {
	slog.Info("Fetching student by Id ")

	stmt, err := s.Db.Prepare("select id, name, email, age from students where id = ? limit 1")
	if err != nil {

		slog.Error("found error while fetching last id ", slog.String("error", err.Error()))

		return types.Student{}, err

	}

	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("no row found in the database")
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}

func (s *Sqlite) GetLists() ([]types.Student, error) {
	slog.Info("fetching all students from data")
	var students []types.Student

	stmt, err := s.Db.Prepare("Select id, name, email, age from students")

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var student types.Student

		err = rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)

		if err != nil {
			return nil, err
		}

		students = append(students, student)

	}

	return students, nil
}
