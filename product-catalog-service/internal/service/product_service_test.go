package service_test

import (
	"errors"
	"testing"

	"github.com/prateek041/product-catalog-service/internal/model"
	"github.com/prateek041/product-catalog-service/internal/service"
)

// MockProductRepository is a mock implementation of the ProductRepository interface for testing.
type MockProductRepository struct {
	products map[string]*model.Product
	err      error
}

func NewMockProductRepository(err error) *MockProductRepository {
	return &MockProductRepository{
		products: make(map[string]*model.Product),
		err:      err,
	}
}

func (m *MockProductRepository) Create(product *model.Product) error {
	if m.err != nil {
		return m.err
	}
	m.products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Get(id string) (*model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	if product, ok := m.products[id]; ok {
		return product, nil
	}
	return nil, nil
}

func (m *MockProductRepository) Update(product *model.Product) error {
	if m.err != nil {
		return m.err
	}
	if _, ok := m.products[product.ID]; ok {
		m.products[product.ID] = product
		return nil
	}
	return errors.New("not found")
}

func (m *MockProductRepository) Delete(id string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.products, id)
	return nil
}

func (m *MockProductRepository) List() ([]*model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	products := make([]*model.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func TestDefaultProductService_CreateProduct(t *testing.T) {
	mockRepo := NewMockProductRepository(nil)
	service := service.NewDefaultProductService(mockRepo)
	product := &model.Product{ID: "test", Name: "Test Product"}

	err := service.CreateProduct(product)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	if _, ok := mockRepo.products["test"]; !ok {
		t.Errorf("Product not created in repository")
	}
}

func TestDefaultProductService_CreateProduct_RepoError(t *testing.T) {
	mockRepo := NewMockProductRepository(errors.New("database error"))
	service := service.NewDefaultProductService(mockRepo)
	product := &model.Product{ID: "test", Name: "Test Product"}

	err := service.CreateProduct(product)
	if err == nil {
		t.Fatalf("CreateProduct should have failed with repository error")
	}
	if err.Error() != "database error" {
		t.Errorf("Incorrect error message: got %q, want %q", err.Error(), "database error")
	}
}
