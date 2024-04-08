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

func (h *ReportHandler) CreateReport(c echo.Context) error {
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
