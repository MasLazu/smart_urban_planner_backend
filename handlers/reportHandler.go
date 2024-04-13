package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"smart_urban_palanner_backend/helper"
	"smart_urban_palanner_backend/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	latitude, err := strconv.ParseFloat(c.FormValue("latitude"), 64)
	if err != nil {
		return helper.NewError(http.StatusBadRequest, "Invalid request body", err)
	}

	longitude, err := strconv.ParseFloat(c.FormValue("longitude"), 64)
	if err != nil {
		return helper.NewError(http.StatusBadRequest, "Invalid request body", err)
	}

	reportRequest := models.ReportRequest{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Address:     c.FormValue("address"),
		Latitude:    latitude,
		Longitude:   longitude,
	}

	if err := c.Validate(reportRequest); err != nil {
		return err
	}

	user := c.Get("user").(jwt.RegisteredClaims)

	image, err := c.FormFile("image")
	if err != nil {
		log.Println(err)
		return helper.NewError(http.StatusBadRequest, "Invalid request body", err)
	}
	src, err := image.Open()
	if err != nil {
		log.Println(err)
		return helper.NewError(http.StatusInternalServerError, "Failed to create report", err)
	}
	defer src.Close()

	id, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		log.Println(err)
		return helper.NewError(http.StatusInternalServerError, "Failed to create report", err)
	}
	image.Filename = id.String() + ".jpg"

	dst, err := os.Create("static/images/" + image.Filename)
	if err != nil {
		log.Println(err)
		return helper.NewError(http.StatusInternalServerError, "Failed to create report", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err)
		return helper.NewError(http.StatusInternalServerError, "Failed to create report", err)
	}

	id, err = uuid.NewV7()
	if err != nil {
		return helper.NewError(http.StatusInternalServerError, "Failed to create report", err)
	}

	report := reportRequest.ToReport()
	report.ID = id.String()
	report.Image = "/static/images/" + image.Filename
	report.Popularity = 0
	report.CreatedAt = time.Now()
	report.AuthorID = user.Subject

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

	user := c.Get("user").(jwt.RegisteredClaims)
	if report.AuthorID != user.Subject {
		return helper.NewError(http.StatusUnauthorized, "Unauthorized", nil)
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

	user := c.Get("user").(jwt.RegisteredClaims)
	if report.AuthorID != user.Subject {
		return helper.NewError(http.StatusUnauthorized, "Unauthorized", nil)
	}

	if err := h.db.Delete(&report).Error; err != nil {
		return helper.NewError(http.StatusInternalServerError, "Failed to delete report", err)
	}

	return c.NoContent(http.StatusNoContent)
}
