package product

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"product-service/api"

	"github.com/google/uuid"
)

var ErrInvalidProduct = errors.New("invalid product")

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	input api.CreateProductRequest,
) (api.Product, error) {
	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		return api.Product{}, fmt.Errorf("%w: name is required", ErrInvalidProduct)
	}

	if input.Variants != nil {
		for i := range *input.Variants {
			variant := &(*input.Variants)[i]

			variant.Name = strings.TrimSpace(variant.Name)
			variant.Sku = strings.TrimSpace(variant.Sku)

			if variant.Name == "" {
				return api.Product{}, fmt.Errorf(
					"%w: variant name is required",
					ErrInvalidProduct,
				)
			}

			if variant.Sku == "" {
				return api.Product{}, fmt.Errorf(
					"%w: variant SKU is required",
					ErrInvalidProduct,
				)
			}

			if variant.Price < 0 {
				return api.Product{}, fmt.Errorf(
					"%w: variant price cannot be negative",
					ErrInvalidProduct,
				)
			}
		}
	}

	return s.repository.Create(ctx, input)
}

func (s *Service) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (api.Product, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) List(
	ctx context.Context,
	params api.ListProductsParams,
) (api.ProductListResponse, error) {
	page := 1
	limit := 20

	if params.Page != nil {
		page = *params.Page
	}

	if params.Limit != nil {
		limit = *params.Limit
	}

	if page < 1 {
		return api.ProductListResponse{}, fmt.Errorf(
			"%w: page must be greater than zero",
			ErrInvalidProduct,
		)
	}

	if limit < 1 || limit > 100 {
		return api.ProductListResponse{}, fmt.Errorf(
			"%w: limit must be between 1 and 100",
			ErrInvalidProduct,
		)
	}

	products, totalItems, err := s.repository.List(ctx, params)
	if err != nil {
		return api.ProductListResponse{}, err
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	return api.ProductListResponse{
		Data: products,
		Pagination: api.Pagination{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}
