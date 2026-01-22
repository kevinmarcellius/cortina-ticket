package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	 quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt) // SIGINT (Ctrl+C)
    signal.Notify(quit, syscall.SIGTERM) // SIGTERM (Kubernetes)
    <-quit // This blocks the main goroutine until a signal is received

    fmt.Println("Shutting down server...")

    // 3. Create a context with a timeout for the shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 4. Attempt graceful shutdown
    if err := e.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    fmt.Println("Server exiting")
}
