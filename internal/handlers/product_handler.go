package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/internal/services"
	"github.com/msyaifudin/pos/internal/validators"
)

type ProductHandler struct {
	ProductService     *services.ProductService
	UserContextService *services.UserContextService
}

func NewProductHandler(productService *services.ProductService, userContextService *services.UserContextService) *ProductHandler {
	return &ProductHandler{ProductService: productService, UserContextService: userContextService}
}

func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	ownerID, err := h.UserContextService.GetOwnerID(userID)
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

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

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

	lang := c.Get("lang").(string)
	if messages := validators.ValidateCreateProduct(product, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

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

	lang := c.Get("lang").(string)
	if messages := validators.ValidateUpdateProduct(product, lang); messages != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": messages,
		})
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

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

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

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

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	products, err := h.ProductService.GetProductsByOutlet(outletUuid, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}
	return JSONSuccess(c, http.StatusOK, "products_retrieved_successfully", products)
}
