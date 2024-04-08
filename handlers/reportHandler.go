package handlers

import (
	"net/http"
	"smart_urban_palanner_backend/helper"
	"smart_urban_palanner_backend/models"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ReportHandler struct {
	db *gorm.DB
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{db: db}
}

func (h *ReportHandler) Create(c echo.Context) error {
	var reportRequest models.ReportRequest
	if err := c.Bind(&reportRequest); err != nil {
		return err
	}

	if err := c.Validate(reportRequest); err != nil {
		return err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return helper.NewError(http.StatusInternalServerError, "Failed to create report", err)
	}

	report := reportRequest.ToReport()
	report.ID = id.String()
	report.Popularity = 0
	report.CreatedAt = time.Now()
	report.AuthorID = "d1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b"
	// report.AuthorID = c.Get("user").(models.User).ID

	if err := h.db.Create(&report).Error; err != nil {
		return helper.NewError(http.StatusInternalServerError, "Failed to create report", err)
	}

	return c.JSON(http.StatusCreated, report)
}

func (h *ReportHandler) Get(c echo.Context) error {
	id := c.Param("id")

	var report models.Report
	if err := h.db.Where("id = ?", id).First(&report).Error; err != nil {
		return helper.NewError(http.StatusNotFound, "Report not found", err)
	}

	return c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) List(c echo.Context) error {
	var reports []models.Report
	if err := h.db.Find(&reports).Error; err != nil {
		return helper.NewError(http.StatusInternalServerError, "Failed to list reports", err)
	}

	return c.JSON(http.StatusOK, reports)
}

func (h *ReportHandler) Update(c echo.Context) error {
	id := c.Param("id")

	var report models.Report
	if err := h.db.Where("id = ?", id).First(&report).Error; err != nil {
		return helper.NewError(http.StatusNotFound, "Report not found", err)
	}

	var reportRequest models.ReportRequest
	if err := c.Bind(&reportRequest); err != nil {
		return err
	}

	if err := c.Validate(reportRequest); err != nil {
		return err
	}

	report.Title = reportRequest.Title
	report.Description = reportRequest.Description
	report.Latitude = reportRequest.Latitude
	report.Longitude = reportRequest.Longitude

	if err := h.db.Save(&report).Error; err != nil {
		return helper.NewError(http.StatusInternalServerError, "Failed to update report", err)
	}

	return c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	var report models.Report
	if err := h.db.Where("id = ?", id).First(&report).Error; err != nil {
		return helper.NewError(http.StatusNotFound, "Report not found", err)
	}

	if err := h.db.Delete(&report).Error; err != nil {
		return helper.NewError(http.StatusInternalServerError, "Failed to delete report", err)
	}

	return c.NoContent(http.StatusNoContent)
}
