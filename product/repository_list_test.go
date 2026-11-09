package product

import (
	"context"
	"testing"

	"product-service/api"
)

func TestRepository_List(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	_, err := repo.Create(ctx, api.CreateProductRequest{
		Name: "Coffee",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(ctx, api.CreateProductRequest{
		Name: "Tea",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Act
	products, totalItems, err := repo.List(
		ctx,
		api.ListProductsParams{},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if totalItems != 2 {
		t.Fatalf("expected 2 total items, got %d", totalItems)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}

	if products[0].Name != "Tea" {
		t.Errorf("expected first product Tea, got %s", products[0].Name)
	}

	if products[1].Name != "Coffee" {
		t.Errorf("expected second product Coffee, got %s", products[1].Name)
	}
}
