package model

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // In a real application, this would be a hashed password
	Email    string `json:"email"`
}
