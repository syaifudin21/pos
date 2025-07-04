package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/services"
)

type ReportHandler struct {
	ReportService *services.ReportService
}

func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{ReportService: reportService}
}

// GetSalesByOutletReport godoc
// @Summary Get sales report by outlet
// @Description Get a sales report for a specific outlet within a date range.
// @Tags Reports
// @Accept json
// @Produce json
// @Param outlet_uuid path string true "Outlet External ID (UUID)"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} SuccessResponse{data=[]models.Order}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/outlets/{outlet_uuid}/sales [get]
func (h *ReportHandler) GetSalesByOutletReport(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Outlet External ID format"})
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid start_date format. Use YYYY-MM-DD"})
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid end_date format. Use YYYY-MM-DD"})
	}

	orders, err := h.ReportService.SalesByOutletReport(outletUuid, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, SuccessResponse{Message: "Sales report by outlet generated successfully", Data: orders})
}

// GetSalesByProductReport godoc
// @Summary Get sales report by product
// @Description Get a sales report for a specific product within a date range.
// @Tags Reports
// @Accept json
// @Produce json
// @Param product_uuid path string true "Product External ID (UUID)"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} SuccessResponse{data=[]models.OrderItem}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reports/products/{product_uuid}/sales [get]
func (h *ReportHandler) GetSalesByProductReport(c echo.Context) error {
	productUuid, err := uuid.Parse(c.Param("product_uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Product External ID format"})
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid start_date format. Use YYYY-MM-DD"})
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid end_date format. Use YYYY-MM-DD"})
	}

	orderItems, err := h.ReportService.SalesByProductReport(productUuid, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, SuccessResponse{Message: "Sales report by product generated successfully", Data: orderItems})
}
