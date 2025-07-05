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

func (h *ReportHandler) GetSalesByOutletReport(c echo.Context) error {
	outletUuid, err := uuid.Parse(c.Param("outlet_uuid"))
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

	orders, err := h.ReportService.SalesByOutletReport(outletUuid, startDate, endDate)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "sales_report_by_outlet_generated_successfully", orders)
}

func (h *ReportHandler) GetSalesByProductReport(c echo.Context) error {
	productUuid, err := uuid.Parse(c.Param("product_uuid"))
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

	orderItems, err := h.ReportService.SalesByProductReport(productUuid, startDate, endDate)
	if err != nil {
		return JSONError(c, MapErrorToStatusCode(err), err.Error())
	}

	return JSONSuccess(c, http.StatusOK, "sales_report_by_product_generated_successfully", orderItems)
}
