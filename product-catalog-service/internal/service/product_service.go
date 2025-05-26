package service

import (
	"github.com/prateek041/product-catalog-service/internal/data"
	"github.com/prateek041/product-catalog-service/internal/model"
)

// ProductService defines the interface for product-related business logic.
type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProduct(id string) (*model.Product, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(id string) error
	ListProducts() ([]*model.Product, error)
}

// DefaultProductService is the default implementation of ProductService.
type DefaultProductService struct {
	repo data.ProductRepository
}

// NewDefaultProductService creates a new DefaultProductService.
func NewDefaultProductService(repo data.ProductRepository) *DefaultProductService {
	return &DefaultProductService{repo: repo}
}

// CreateProduct implements the business logic for creating a product.
func (s *DefaultProductService) CreateProduct(product *model.Product) error {
	// Add any business logic/validation here before saving
	return s.repo.Create(product)
}

// GetProduct implements the business logic for retrieving a product.
func (s *DefaultProductService) GetProduct(id string) (*model.Product, error) {
	return s.repo.Get(id)
}

// UpdateProduct implements the business logic for updating a product.
func (s *DefaultProductService) UpdateProduct(product *model.Product) error {
	// Add any business logic/validation here before updating
	return s.repo.Update(product)
}

// DeleteProduct implements the business logic for deleting a product.
func (s *DefaultProductService) DeleteProduct(id string) error {
	// Add any business logic here before deleting
	return s.repo.Delete(id)
}

// ListProducts implements the business logic for listing all products.
func (s *DefaultProductService) ListProducts() ([]*model.Product, error) {
	return s.repo.List()
}
