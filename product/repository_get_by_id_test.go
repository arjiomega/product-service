package product

import (
	"context"
	"errors"
	"product-service/api"
	"testing"

	"github.com/google/uuid"
)

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
