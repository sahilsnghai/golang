package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sahilsnghai/Project6/types"
)

type mockUserStore struct{}

func (s *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if slices.Contains([]string{"sahil1@g.com"}, email) {
		return &types.User{Email: email}, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *mockUserStore) GetUserById(id int) (*types.User, error) {
	return nil, nil
}

func (s *mockUserStore) CreateUser(user types.User) error {
	return nil
}

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := Newhandler(userStore)

	t.Run("should fail if the user payload invalad",
		func(t *testing.T) {
			payload := types.RegisterUserPayload{
				FirstName: "asjd",
				LastName:  "12",
				Email:     "invalid",
				Password:  "Asda",
			}

			marshedPayload, _ := json.Marshal(payload)
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshedPayload))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/register", handler.handleRegister)
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
			}
		})

	t.Run("should pass if the user payload is valid",
		func(t *testing.T) {
			payload := types.RegisterUserPayload{
				FirstName: "asjd",
				LastName:  "12",
				Email:     "valid@gmail.com",
				Password:  "Asda",
			}

			marshedPayload, _ := json.Marshal(payload)
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshedPayload))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/register", handler.handleRegister)
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
			}
		})

}

func TestUserLogineHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := Newhandler(userStore)

	t.Run("should fail if the user email or password is invalid",
		func(t *testing.T) {
			payload := types.LoginUserPayload{
				Email:    "sahil@g.com",
				Password: "abcd",
			}

			marshedPayload, _ := json.Marshal(payload)
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshedPayload))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/login", handler.handleLogin)
			router.ServeHTTP(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
			}
		})

	t.Run("should pass if the user email and password is valid",
		func(t *testing.T) {
			payload := types.LoginUserPayload{
				Email:    "sahil1@g.com",
				Password: "abc",
			}

			marshedPayload, _ := json.Marshal(payload)
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshedPayload))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/login", handler.handleLogin)
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
				t.Errorf("payload response %s", rr.Body)
			}
		})

}
