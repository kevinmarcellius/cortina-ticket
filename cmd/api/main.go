package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kevinmarcellius/cortina-ticket/config"
	"github.com/kevinmarcellius/cortina-ticket/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Health check endpoints
	healthHandler := handler.NewHealthHandler(cfg.DB)
	e.GET("/health/live", healthHandler.Live)
	e.GET("/health/ready", healthHandler.Ready)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
