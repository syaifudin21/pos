package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/services"
)

func isValidProductType(productType string) bool {
	for _, pt := range models.AllowedProductTypes {
		if pt == productType {
			return true
		}
	}
	return false
}

type ProductHandler struct {
	ProductService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get a list of all products.
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse{data=[]models.Product}
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.ProductService.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Products retrieved successfully", Data: products})
}

// GetProductByID godoc
// @Summary Get product by External ID
// @Description Get a single product by its External ID.
// @Tags Products
// @Accept json
// @Produce json
// @Param external_id path string true "Product External ID (UUID)"
// @Success 200 {object} SuccessResponse{data=models.Product}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{external_id} [get]
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("external_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid External ID format"})
	}
	product, err := h.ProductService.GetProductByUuid(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Product retrieved successfully", Data: product})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details.
// @Tags Products
// @Accept json
// @Produce json
// @Param product body ProductCreateRequest true "Product details"
// @Success 201 {object} SuccessResponse{data=models.Product}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	product := new(ProductCreateRequest)
	if err := c.Bind(product); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	if !isValidProductType(product.Type) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid product type specified"})
	}

	newProduct := &models.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Type:        product.Type,
	}

	createdProduct, err := h.ProductService.CreateProduct(newProduct)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, SuccessResponse{Message: "Product created successfully", Data: createdProduct})
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update an existing product by its External ID.
// @Tags Products
// @Accept json
// @Produce json
// @Param external_id path string true "Product External ID (UUID)"
// @Param product body ProductUpdateRequest true "Updated product details"
// @Success 200 {object} SuccessResponse{data=models.Product}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{external_id} [put]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("external_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid External ID format"})
	}
	product := new(ProductUpdateRequest)
	if err := c.Bind(product); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request payload"})
	}

	if !isValidProductType(product.Type) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid product type specified"})
	}

	updatedProduct := &models.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		Type:        product.Type,
	}

	result, err := h.ProductService.UpdateProduct(id, updatedProduct)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, SuccessResponse{Message: "Product updated successfully", Data: result})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its Uuid.
// @Tags Products
// @Accept json
// @Produce json
// @Param uuid path string true "Product Uuid"
// @Success 204 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{uuid} [delete]
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid UUID format"})
	}
	err = h.ProductService.DeleteProduct(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNoContent, SuccessResponse{Message: "Product deleted successfully"})
}

type ProductCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku,omitempty"`
	Type        string  `json:"type"`
}

type ProductUpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku,omitempty"`
	Type        string  `json:"type"`
}
