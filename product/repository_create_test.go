package product

import (
	"context"
	"product-service/api"
	"testing"

	"github.com/google/uuid"
)

func TestRepository_CreateAndGetByID(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	variants := []api.CreateProductVariantRequest{
		{
			Name:  "Small",
			Sku:   "PRODUCT-SMALL",
			Price: 100.00,
		},
		{
			Name:  "Large",
			Sku:   "PRODUCT-LARGE",
			Price: 150.00,
		},
	}

	description := "Product Description"

	input := api.CreateProductRequest{
		Name:        "Product Name",
		Description: &description,
		Variants:    &variants,
	}

	// Act: Create product
	createdProduct, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	// Act: Retrieve product
	retrievedProduct, err := repo.GetByID(ctx, createdProduct.Id)
	if err != nil {
		t.Fatalf("failed to get product by ID: %v", err)
	}

	// Assert: Product fields
	if retrievedProduct.Id != createdProduct.Id {
		t.Errorf("expected product ID %v, got %v",
			createdProduct.Id, retrievedProduct.Id)
	}

	if retrievedProduct.Name != input.Name {
		t.Errorf("expected product name %q, got %q",
			input.Name, retrievedProduct.Name)
	}

	if retrievedProduct.Description == nil ||
		*retrievedProduct.Description != *input.Description {
		t.Errorf("expected description %q, got %v",
			*input.Description, retrievedProduct.Description)
	}

	if retrievedProduct.Status != api.Draft {
		t.Errorf("expected status %q, got %q",
			api.Draft, retrievedProduct.Status)
	}

	// Assert: Variants
	if retrievedProduct.Variants == nil {
		t.Fatal("expected variants, got nil")
	}

	if len(*retrievedProduct.Variants) != len(variants) {
		t.Fatalf("expected %d variants, got %d",
			len(variants), len(*retrievedProduct.Variants))
	}

	for i, expected := range variants {
		actual := (*retrievedProduct.Variants)[i]

		if actual.Name != expected.Name {
			t.Errorf("variant %d: expected name %q, got %q",
				i, expected.Name, actual.Name)
		}

		if actual.Sku != expected.Sku {
			t.Errorf("variant %d: expected SKU %q, got %q",
				i, expected.Sku, actual.Sku)
		}

		if actual.Price != expected.Price {
			t.Errorf("variant %d: expected price %.2f, got %.2f",
				i, expected.Price, actual.Price)
		}

		if actual.Status != api.Inactive {
			t.Errorf("variant %d: expected status %q, got %q",
				i, api.Inactive, actual.Status)
		}

		if actual.Id == uuid.Nil {
			t.Errorf("variant %d: expected generated ID, got empty ID", i)
		}
	}

}
