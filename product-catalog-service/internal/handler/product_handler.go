package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prateek041/product-catalog-service/internal/model"
	"github.com/prateek041/product-catalog-service/internal/service"
)

// ProductHandler handles HTTP requests related to products.
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct handles the creation of a new product.
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.service.CreateProduct(&product)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product) // Optionally return the created product
}

// GetProduct handles retrieving a product by ID.
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.service.GetProduct(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving product: %v", err), http.StatusInternalServerError)
		return
	}
	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct handles updating an existing product.
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	product.ID = id // Ensure the ID from the path matches the product to be updated

	err = h.service.UpdateProduct(&product)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product) // Optionally return the updated product
}

// DeleteProduct handles deleting a product by ID.
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteProduct(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content to return for successful deletion
}

// ListProducts handles retrieving a list of all products.
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing products: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
