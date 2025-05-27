package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prateek041/user-management-service/internal/data"
	"github.com/prateek041/user-management-service/internal/handler"
	"github.com/prateek041/user-management-service/internal/service"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}

func main() {
	userRepo := data.NewInMemoryUserRepository()
	userService := service.NewDefaultUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)

	r := mux.NewRouter()

	r.HandleFunc("/health", healthHandler) // health check handler.

	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	port := ":8081"
	log.Printf("User Management Service listening on port %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
