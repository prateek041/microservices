package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prateek041/product-catalog-service/internal/data"
	"github.com/prateek041/product-catalog-service/internal/handler"
	"github.com/prateek041/product-catalog-service/internal/service"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}

func main() {
	// Initialize the in-memory product repository
	productRepo := data.NewInMemoryProductRepository()

	// Initialize the product service, passing in the repository
	productService := service.NewDefaultProductService(productRepo)

	// Initialize the product handler, passing in the service
	productHandler := handler.NewProductHandler(productService)

	// Create a new router using gorilla/mux
	r := mux.NewRouter()

	// Define API routes and associate them with handler functions
	r.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/products", productHandler.ListProducts).Methods("GET")

	// Add the health check endpoint.
	r.HandleFunc("/health", healthHandler).Methods("GET")

	// Start the HTTP server
	port := ":8080" // Can be configured via environment variables later on.
	log.Printf("Product Catalog Service listening on port %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
