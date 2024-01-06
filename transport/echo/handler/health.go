package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status,omitempty"`
}

func HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{Status: "OK"})
}
