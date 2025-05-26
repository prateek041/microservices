package data

import (
	"fmt"
	"sync"

	"github.com/prateek041/product-catalog-service/internal/model"
)

// ProductRepository defines the interface for interacting with product data.
type ProductRepository interface {
	Create(product *model.Product) error
	Get(id string) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id string) error
	List() ([]*model.Product, error)
}

// InMemoryProductRepository is an in-memory implementation of InMemoryProductRepository.
type InMemoryProductRepository struct {
	products map[string]*model.Product
	mu       sync.RWMutex
}

// NewInMemoryProductRepository creates a new InMemoryProductRepository.
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*model.Product),
	}
}

// Create adds a new product to the in-memory store.
func (r *InMemoryProductRepository) Create(product *model.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; exists {
		return fmt.Errorf("product with ID %s already exists", product.ID)
	}

	r.products[product.ID] = product
	return nil
}

// Get retrieves a product by its ID from the data store.
func (r *InMemoryProductRepository) Get(id string) (*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	product, exists := r.products[id]
	if !exists {
		return nil, fmt.Errorf("product with ID %s not found", id)
	}
	return product, nil
}

// Update updates an existing product in the store.

func (r *InMemoryProductRepository) Update(product *model.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return fmt.Errorf("product with Id %s not found", product.ID)
	}

	r.products[product.ID] = product
	return nil
}

// Delete removes a product form the store.
func (r *InMemoryProductRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.products, id)
	return nil
}

// List retrieves all products from the store.
func (r *InMemoryProductRepository) List() ([]*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*model.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}

	return products, nil
}
