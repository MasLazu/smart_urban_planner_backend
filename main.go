package main

import (
	"net/http"
	"smart_urban_palanner_backend/handlers"
	"smart_urban_palanner_backend/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	logger "github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=smart_urban_planner port=5432 sslmode=disable TimeZone=Asia/Shanghai"))
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Report{}, &models.User{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(logger.Logger())

	reportHandler := handlers.NewReportHandler(db)

	e.POST("/reports", reportHandler.Create)

	e.Logger.Fatal(e.Start(":8080"))
}
