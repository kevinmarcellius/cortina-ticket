package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// HealthHandler handles health checks
type HealthHandler struct {
	DB *gorm.DB
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

// Live is a liveness probe
func (h *HealthHandler) Live(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "live"})
}

// Ready is a readiness probe
func (h *HealthHandler) Ready(c echo.Context) error {
	sqlDB, err := h.DB.DB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "not ready"})
	}

	err = sqlDB.Ping()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "not ready"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
}
