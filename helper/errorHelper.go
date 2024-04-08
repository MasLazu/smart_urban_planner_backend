package helper

import "github.com/labstack/echo/v4"

func NewError(code int, message string, error error) *echo.HTTPError {
	return &echo.HTTPError{Message: message, Code: code, Internal: error}
}
