package data_test

import (
	"reflect"
	"testing"

	"github.com/prateek041/product-catalog-service/internal/data"
	"github.com/prateek041/product-catalog-service/internal/model"
)

func TestInMemoryProductRepository_Create_Get(t *testing.T) {
	repo := data.NewInMemoryProductRepository()
	product := &model.Product{ID: "1", Name: "Test Product"}

	err := repo.Create(product)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrievedProduct, err := repo.Get("1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !reflect.DeepEqual(product, retrievedProduct) {
		t.Errorf("Retrieved product is not the same as created product: got %+v, want %+v", retrievedProduct, product)
	}
}

func TestInMemoryProductRepository_Get_NotFound(t *testing.T) {
	repo := data.NewInMemoryProductRepository()

	_, err := repo.Get("nonexistent")
	if err == nil {
		t.Fatalf("Get should have failed for non-existent ID")
	}
}

func TestInMemoryProductRepository_Update(t *testing.T) {
	repo := data.NewInMemoryProductRepository()
	product := &model.Product{ID: "2", Name: "Initial Name", Price: 10.0}
	repo.Create(product)

	updatedProduct := &model.Product{ID: "2", Name: "Updated Name", Price: 20.0}
	err := repo.Update(updatedProduct)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	retrievedProduct, _ := repo.Get("2")
	if !reflect.DeepEqual(updatedProduct, retrievedProduct) {
		t.Errorf("Retrieved product after update is incorrect: got %+v, want %+v", retrievedProduct, updatedProduct)
	}
}

func TestInMemoryProductRepository_Delete(t *testing.T) {
	repo := data.NewInMemoryProductRepository()
	product := &model.Product{ID: "3", Name: "Product to Delete"}
	repo.Create(product)

	err := repo.Delete("3")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.Get("3")
	if err == nil {
		t.Fatalf("Get after delete should have failed")
	}
}

func TestInMemoryProductRepository_List(t *testing.T) {
	repo := data.NewInMemoryProductRepository()
	product1 := &model.Product{ID: "4", Name: "Product A"}
	product2 := &model.Product{ID: "5", Name: "Product B"}
	repo.Create(product1)
	repo.Create(product2)

	products, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(products) != 2 {
		t.Errorf("List should return 2 products, got %d", len(products))
	}

}
