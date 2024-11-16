package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sahilsnghai/Project6/config"
	"github.com/sahilsnghai/Project6/service/auth"
	"github.com/sahilsnghai/Project6/types"
	"github.com/sahilsnghai/Project6/utils"
)

type Handler struct {
	store types.UserStore
}

func Newhandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(route *mux.Router) {
	route.HandleFunc("/login", h.handleLogin).Methods("POST")
	route.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// payload validation
	if err := utils.Validate.Struct(payload); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	// password hashing
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}
	secret := []byte(config.EnvsInit.SecretKey)
	token, err := auth.CreateJWT(secret, u.Id)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token})

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("Registing New user")
	var payload types.RegisterUserPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// if user exists
	if _, err := h.store.GetUserByEmail(payload.Email); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	// payload validation
	if err := utils.Validate.Struct(payload); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errs))
		return
	}

	// password hashing
	hashedPassword, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		log.Fatal(err)
		utils.WriteError(w, http.StatusInternalServerError, err)

	}
	utils.WriteJson(w, http.StatusCreated, nil)
}
