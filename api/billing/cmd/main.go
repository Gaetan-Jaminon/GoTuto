package main

import (
	"fmt"
	"log"
	"strings"
	
	"gotuto/api/billing/internal/config"
	"gotuto/api/billing/internal/database"
	"gotuto/api/billing/internal/handlers"
	
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()
	cfg.Print() // Log configuration (without sensitive data)
	
	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	// Set up router
	router := setupRouter(cfg)
	
	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s in %s mode", addr, cfg.Server.Mode)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(cfg *config.Config) *gin.Engine {
	// Set Gin mode based on config
	gin.SetMode(cfg.Server.Mode)
	
	router := gin.Default()
	
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// CORS middleware from config
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range cfg.CORS.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.CORS.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.CORS.AllowedHeaders, ", "))
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "demo01-api",
		})
	})
	
	// API routes
	api := router.Group("/api/v1")
	{
		// Client routes
		clients := api.Group("/clients")
		{
			clients.GET("", handlers.GetClients)
			clients.GET("/:id", handlers.GetClient)
			clients.POST("", handlers.CreateClient)
			clients.PUT("/:id", handlers.UpdateClient)
			clients.DELETE("/:id", handlers.DeleteClient)
			clients.GET("/:client_id/invoices", handlers.GetInvoicesByClient)
		}
		
		// Invoice routes
		invoices := api.Group("/invoices")
		{
			invoices.GET("", handlers.GetInvoices)
			invoices.GET("/:id", handlers.GetInvoice)
			invoices.POST("", handlers.CreateInvoice)
			invoices.PUT("/:id", handlers.UpdateInvoice)
			invoices.DELETE("/:id", handlers.DeleteInvoice)
		}
	}
	
	return router
}