package main

import (
	"net/http"
	"smart_urban_palanner_backend/handlers"
	"smart_urban_palanner_backend/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
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

func authMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			exp := int64(claims["exp"].(float64))

			c.Set("user", jwt.RegisteredClaims{
				Subject:   claims["sub"].(string),
				ExpiresAt: jwt.NewNumericDate(time.Unix(exp, 0)),
			})
		},
	})
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
	authHandler := handlers.NewAuthHandler(db)

	auth := e.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)

	report := e.Group("/reports")
	report.GET("/:id", reportHandler.Get)
	report.GET("", reportHandler.List)
	report.POST("", reportHandler.Create, authMiddleware())
	report.PUT("/:id", reportHandler.Update, authMiddleware())
	report.DELETE("/:id", reportHandler.Delete, authMiddleware())

	e.Logger.Fatal(e.Start(":8080"))
}
