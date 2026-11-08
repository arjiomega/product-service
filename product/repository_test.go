package product

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"product-service/api"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

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

func TestRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	nonExistentID := api.ID(uuid.New())

	// Act
	_, err := repo.GetByID(ctx, nonExistentID)

	// Assert
	if !errors.Is(err, ErrProductNotFound) {
		t.Errorf(
			"expected ErrProductNotFound, got %v",
			err,
		)
	}
}

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

	if products[0].Name != "Coffee" {
		t.Errorf("expected first product Coffee, got %s", products[0].Name)
	}

	if products[1].Name != "Tea" {
		t.Errorf("expected second product Tea, got %s", products[1].Name)
	}
}

func TestRepository_Create_CancelledContext(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	input := api.CreateProductRequest{
		Name: "Coffee",
	}

	// Act
	_, err := repo.Create(ctx, input)

	// Assert
	if !errors.Is(err, context.Canceled) {
		t.Errorf(
			"expected context.Canceled, got %v",
			err,
		)
	}
}

func TestRepository_GetByID_CancelledContext(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cancel()

	nonExistentID := api.ID(uuid.New())

	// Act
	_, err := repo.GetByID(ctx, nonExistentID)

	// Assert
	if !errors.Is(err, context.Canceled) {
		t.Errorf(
			"expected context.Canceled, got %v",
			err,
		)
	}
}
