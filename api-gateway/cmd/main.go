package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"api-gateway/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config := loadConfig()

	// Initialize Gin router
	r := gin.Default()

	// Apply global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Initialize handlers
	membershipHandler := handler.NewMembershipHandler(config.MembershipServiceURL)
	websiteHandler := handler.NewWebsiteHandler(config.WebsiteServiceURL)
	bankingHandler := handler.NewBankingHandler(config.BankingServiceURL)

	// Setup routes
	router.SetupMembershipRoutes(r, membershipHandler)
	router.SetupWebsiteRoutes(r, websiteHandler)
	router.SetupBankingRoutes(r, bankingHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:    config.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("API Gateway starting on %s", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

type Config struct {
	Port                  string
	MembershipServiceURL  string
	WebsiteServiceURL     string
	BankingServiceURL     string
}

func loadConfig() *Config {
	return &Config{
		Port:                  getEnv("PORT", ":8080"),
		MembershipServiceURL:  getEnv("MEMBERSHIP_SERVICE_URL", "http://localhost:8081"),
		WebsiteServiceURL:     getEnv("WEBSITE_SERVICE_URL", "http://localhost:8082"),
		BankingServiceURL:     getEnv("BANKING_SERVICE_URL", "http://localhost:8083"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

