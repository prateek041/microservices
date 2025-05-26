package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prateek041/user-management-service/internal/data"
	"github.com/prateek041/user-management-service/internal/handler"
	"github.com/prateek041/user-management-service/internal/service"
)

func main() {
	userRepo := data.NewInMemoryUserRepository()
	userService := service.NewDefaultUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/auth/login", userHandler.AuthenticateUser).Methods("POST")

	port := ":8081"
	log.Printf("User Management Service listening on port %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
