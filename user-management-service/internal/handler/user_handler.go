package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prateek041/user-management-service/internal/model"
	"github.com/prateek041/user-management-service/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.HashedPassword = string(hashedPassword)
	user.Password = "" // We don't save the plain text password.

	err = h.service.CreateUser(&user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// func (h *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
// 	var credentials struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}
// 	err := json.NewDecoder(r.Body).Decode(&credentials)
// 	if err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}
//
// 	user, err := h.service.AuthenticateUser(credentials.Username, credentials.Password)
// 	if err != nil {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}
//
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Authentication successful", "user": user})
// }
