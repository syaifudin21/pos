package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
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
// @Success 200 {object} SuccessResponse{data=[]dtos.ProductResponse}
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.ProductService.GetAllProducts()
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	var productResponses []dtos.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, dtos.ProductResponse{
			ID:          product.ID,
			Uuid:        product.Uuid,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			SKU:         product.SKU,
			Type:        product.Type,
		})
	}
	return JSONSuccess(c, http.StatusOK, "products_retrieved_successfully", productResponses)
}

// GetProductByID godoc
// @Summary Get product by External ID
// @Description Get a single product by its External ID.
// @Tags Products
// @Accept json
// @Produce json
// @Param external_id path string true "Product External ID (UUID)"
// @Success 200 {object} SuccessResponse{data=dtos.ProductResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{external_id} [get]
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("external_id"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_external_id_format")
	}
	product, err := h.ProductService.GetProductByUuid(id)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "product_retrieved_successfully", product)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details.
// @Tags Products
// @Accept json
// @Produce json
// @Param product body dtos.ProductCreateRequest true "Product details"
// @Success 201 {object} SuccessResponse{data=dtos.ProductResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	product := new(dtos.ProductCreateRequest)
	if err := c.Bind(product); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	if !isValidProductType(product.Type) {
		return JSONError(c, http.StatusBadRequest, "invalid_product_type_specified")
	}

	createdProduct, err := h.ProductService.CreateProduct(product)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusCreated, "product_created_successfully", createdProduct)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update an existing product by its External ID.
// @Tags Products
// @Accept json
// @Produce json
// @Param uuid path string true "Product External ID (UUID)"
// @Param product body dtos.ProductUpdateRequest true "Updated product details"
// @Success 200 {object} SuccessResponse{data=dtos.ProductResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{uuid} [put]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_external_id_format")
	}
	product := new(dtos.ProductUpdateRequest)
	if err := c.Bind(product); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	if !isValidProductType(product.Type) {
		return JSONError(c, http.StatusBadRequest, "invalid_product_type_specified")
	}

	result, err := h.ProductService.UpdateProduct(id, product)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "product_updated_successfully", result)
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
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	err = h.ProductService.DeleteProduct(id)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "product_deleted_successfully", nil)
}

// GetProductsByOutlet godoc
// @Summary Get products by outlet
// @Description Get a list of products available in a specific outlet.
// @Tags Products
// @Accept json
// @Produce json
// @Param outlet_uuid path string true "Outlet Uuid"
// @Success 200 {object} SuccessResponse{data=[]dtos.ProductOutletResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /outlets/{outlet_uuid}/products [get]
func (h *ProductHandler) GetProductsByOutlet(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}

	products, err := h.ProductService.GetProductsByOutlet(outletUuid)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "products_retrieved_successfully", products)
}
