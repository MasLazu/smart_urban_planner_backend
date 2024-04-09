package handlers

import (
	"net/http"
	"smart_urban_palanner_backend/helper"
	"smart_urban_palanner_backend/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}
func (h *AuthHandler) Login(c echo.Context) error {
	var userRequest models.UserLogin
	if err := c.Bind(&userRequest); err != nil {
		return helper.NewError(http.StatusBadRequest, "Invalid request body", err)
	}

	if err := c.Validate(userRequest); err != nil {
		return helper.NewError(http.StatusBadRequest, "Invalid request body", err)
	}

	user := userRequest.ToUser()
	if err := h.db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		return helper.NewError(http.StatusUnauthorized, "Invalid email or password", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password)); err != nil {
		return helper.NewError(http.StatusUnauthorized, "Invalid email or password", err)
	}

	claims := jwt.RegisteredClaims{
		Subject:   user.ID,
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

func (h *AuthHandler) Register(c echo.Context) error {
	var userRequest models.UserRegister
	if err := c.Bind(&userRequest); err != nil {
		return helper.NewError(http.StatusBadRequest, "Invalid request body", err)
	}

	if err := c.Validate(userRequest); err != nil {
		return helper.NewError(http.StatusBadRequest, "Invalid request body", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return helper.NewError(http.StatusInternalServerError, "Failed to create report", err)
	}

	user := userRequest.ToUser()
	user.ID = id.String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return helper.NewError(http.StatusInternalServerError, "Internal server error", err)
	}

	user.Password = string(hashedPassword)

	if err := h.db.Create(&user).Error; err != nil {
		return helper.NewError(http.StatusInternalServerError, "Email already taken", err)
	}

	return c.JSON(http.StatusCreated, user)
}
