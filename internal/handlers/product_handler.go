package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/pkg/utils"
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

func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	ownerID, err := h.ProductService.GetOwnerID(userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	products, err := h.ProductService.GetAllProducts(ownerID)
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

func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	product, err := h.ProductService.GetProductByUuid(id, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "product_retrieved_successfully", product)
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	product := new(dtos.ProductCreateRequest)
	if err := c.Bind(product); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	if !isValidProductType(product.Type) {
		return JSONError(c, http.StatusBadRequest, "invalid_product_type_specified")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	createdProduct, err := h.ProductService.CreateProduct(product, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusCreated, "product_created_successfully", createdProduct)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}
	product := new(dtos.ProductUpdateRequest)
	if err := c.Bind(product); err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_request_payload")
	}

	if !isValidProductType(product.Type) {
		return JSONError(c, http.StatusBadRequest, "invalid_product_type_specified")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	result, err := h.ProductService.UpdateProduct(id, product, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "product_updated_successfully", result)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_uuid_format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	err = h.ProductService.DeleteProduct(id, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusNoContent, "product_deleted_successfully", nil)
}

func (h *ProductHandler) GetProductsByOutlet(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*utils.Claims)
	userID := claims.ID

	products, err := h.ProductService.GetProductsByOutlet(outletUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "products_retrieved_successfully", products)
}
