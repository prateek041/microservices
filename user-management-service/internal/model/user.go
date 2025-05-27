package model

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Password       string `json:"password"` // In a real application, this would be a hashed password
	HashedPassword string `json:"hashedPassword"`
	Email          string `json:"email"`
}
