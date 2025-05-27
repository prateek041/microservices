package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prateek041/user-management-service/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	userService service.UserService
	signingKey  []byte
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	signingKey := []byte(os.Getenv("JWT_SIGNING_KEY")) // Get the sercret key to sign JWTs.

	if len(signingKey) == 0 {
		// For production, we will be using more secure method for signing JWTs.
		signingKey = []byte("sample-signing-key")
		log.Printf("Warning: Using insecure default JWT signing Key!")
	}

	return &AuthHandler{
		userService: userService,
		signingKey:  signingKey,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByUsername(req.Username)
	if err != nil {
		log.Println("get user failed")
		log.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		log.Println("password comparison failed")
		log.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := jwt.RegisteredClaims{
		Subject:   user.ID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		// we can add more claims here like roles etc.
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Using HMAC for now, we'll switch to RSA later.
	tokenString, err := token.SignedString(h.signingKey)

	if err != nil {
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	resp := LoginResponse{Token: tokenString}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
