package product

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"product-service/api"

	"github.com/google/uuid"
)

var ErrProductNotFound = errors.New("product not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	input api.CreateProductRequest,
) (api.Product, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return api.Product{}, err
	}
	defer tx.Rollback()

	now := time.Now().UTC()

	product := api.Product{
		Name:        input.Name,
		Description: input.Description,
		Status:      api.Draft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = tx.QueryRowContext(
		ctx,
		`
		INSERT INTO products (
			name, description, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`,
		product.Name,
		product.Description,
		product.Status,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.Id)

	if err != nil {
		return api.Product{}, err
	}

	variants := make([]api.ProductVariant, 0)

	if input.Variants != nil {
		for _, variant := range *input.Variants {
			productVariant := api.ProductVariant{
				Name:      variant.Name,
				Sku:       variant.Sku,
				Price:     variant.Price,
				Status:    api.Inactive,
				CreatedAt: now,
				UpdatedAt: now,
			}

			err := tx.QueryRowContext(
				ctx,
				`
				INSERT INTO product_variants (
					product_id, name, sku, price,
					status, created_at, updated_at
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id
				`,
				product.Id,
				productVariant.Name,
				productVariant.Sku,
				productVariant.Price,
				productVariant.Status,
				productVariant.CreatedAt,
				productVariant.UpdatedAt,
			).Scan(&productVariant.Id)

			if err != nil {
				return api.Product{}, err
			}

			variants = append(variants, productVariant)
		}
	}

	product.Variants = &variants

	if err := tx.Commit(); err != nil {
		return api.Product{}, err
	}

	return product, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (api.Product, error) {
	var product api.Product

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, name, description, status, created_at, updated_at
		FROM products
		WHERE id = $1
		`,
		id,
	).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return api.Product{}, ErrProductNotFound
	}
	if err != nil {
		return api.Product{}, err
	}

	variants, err := r.getVariants(ctx, id)
	if err != nil {
		return api.Product{}, err
	}

	product.Variants = &variants

	return product, nil
}

func (r *Repository) getVariants(
	ctx context.Context,
	productID uuid.UUID,
) ([]api.ProductVariant, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT id, name, sku, price, status, created_at, updated_at
		FROM product_variants
		WHERE product_id = $1
		ORDER BY created_at ASC
		`,
		productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []api.ProductVariant

	for rows.Next() {
		var variant api.ProductVariant

		if err := rows.Scan(
			&variant.Id,
			&variant.Name,
			&variant.Sku,
			&variant.Price,
			&variant.Status,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		); err != nil {
			return nil, err
		}

		variants = append(variants, variant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if variants == nil {
		variants = make([]api.ProductVariant, 0)
	}

	return variants, nil
}

func (r *Repository) List(
	ctx context.Context,
	params api.ListProductsParams,
) ([]api.Product, int, error) {
	page := 1
	limit := 20

	if params.Page != nil {
		page = *params.Page
	}

	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := (page - 1) * limit

	var status string
	var search string

	if params.Status != nil {
		status = string(*params.Status)
	}

	if params.Search != nil {
		search = *params.Search
	}

	var totalItems int

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT COUNT(*)
		FROM products
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR name ILIKE '%' || $2 || '%')
		`,
		status,
		search,
	).Scan(&totalItems)

	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT id, name, description, status, created_at, updated_at
		FROM products
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR name ILIKE '%' || $2 || '%')
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
		`,
		status,
		search,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]api.Product, 0)

	for rows.Next() {
		var product api.Product

		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}
