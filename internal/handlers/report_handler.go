package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/services"
)

type ReportHandler struct {
	ReportService      *services.ReportService
	UserContextService *services.UserContextService
}

func NewReportHandler(reportService *services.ReportService, userContextService *services.UserContextService) *ReportHandler {
	return &ReportHandler{ReportService: reportService, UserContextService: userContextService}
}

func (h *ReportHandler) GetSalesByOutletReport(c echo.Context) error {
	outletUuidParam := c.Param("outlet_uuid")
	outletUuid, err := uuid.Parse(outletUuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_outlet_uuid_format")
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_start_date_format")
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_end_date_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	orders, err := h.ReportService.SalesByOutletReport(outletUuid, startDate, endDate, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "sales_report_by_outlet_generated_successfully", orders)
}

func (h *ReportHandler) GetSalesByProductReport(c echo.Context) error {
	productUuidParam := c.Param("product_uuid")
	productUuid, err := uuid.Parse(productUuidParam)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_product_uuid_format")
	}

	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_start_date_format")
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return JSONError(c, http.StatusBadRequest, "invalid_end_date_format")
	}

	userID, err := h.UserContextService.GetUserIDFromEchoContext(c)
	if err != nil {
		return JSONError(c, http.StatusUnauthorized, err.Error())
	}

	orderItems, err := h.ReportService.SalesByProductReport(productUuid, startDate, endDate, userID)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "sales_report_by_product_generated_successfully", orderItems)
}
