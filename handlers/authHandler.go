package handlers

import (
	"net/http"
	"smart_urban_palanner_backend/helper"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}
func (h *AuthHandler) Login(c echo.Context) error {
	claims := jwt.RegisteredClaims{
		Subject:   "d1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return helper.NewError(http.StatusInternalServerError, "Internal server error", err)
	}
	signedToken = "Bearer " + signedToken

	c.Response().Header().Set("Authorization", signedToken)

	return c.JSON(http.StatusOK, map[string]string{
		"token": signedToken,
	})
}
