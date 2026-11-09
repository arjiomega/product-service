package product

import (
	"context"
	"errors"
	"product-service/api"
	"testing"

	"github.com/google/uuid"
)

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
