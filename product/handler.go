package product

import (
	"errors"
	"net/http"

	"product-service/api"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// POST /products
func (h *Handler) CreateProduct(c *gin.Context) {
	var input api.CreateProductRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{
			Error: "Invalid request body",
		})
		return
	}

	product, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrInvalidProduct) {
			c.JSON(http.StatusBadRequest, api.Error{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, api.Error{
			Error: "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GET /products/:product_id
func (h *Handler) GetProduct(c *gin.Context, productID api.ID) {
	id := uuid.UUID(productID)

	product, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.JSON(http.StatusNotFound, api.Error{
				Error: "Product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, api.Error{
			Error: "Failed to retrieve product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GET /products
func (h *Handler) ListProducts(
	c *gin.Context,
	params api.ListProductsParams,
) {
	result, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		if errors.Is(err, ErrInvalidProduct) {
			c.JSON(http.StatusBadRequest, api.Error{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, api.Error{
			Error: "Failed to retrieve products",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Ensure Handler implements the generated interface.
var _ api.ServerInterface = (*Handler)(nil)
