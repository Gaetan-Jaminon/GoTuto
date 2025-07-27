package main

import (
	"fmt"
	"log"
	"strings"
	
	"gaetanjaminon/GoTuto/internal/billing/config"
	"gaetanjaminon/GoTuto/internal/billing/database"
	"gaetanjaminon/GoTuto/internal/billing/api"
	
	"github.com/gin-gonic/gin"
)

func main() {
	// Load billing domain configuration
	cfg := config.MustLoad()
	
	log.Println("=== Billing Service Configuration ===")
	log.Printf("Server: Port=%d, Mode=%s", cfg.Server.Port, cfg.Server.Mode)
	log.Printf("Database: Host=%s:%d, Name=%s, Schema=%s, User=%s", 
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.Schema, cfg.Database.Username)
	log.Printf("Logging: Level=%s, Format=%s", cfg.Logging.Level, cfg.Logging.Format)
	
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

func setupRouter(cfg *config.BillingConfig) *gin.Engine {
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
			"service": "billing-api",
			"domain":  "billing",
		})
	})
	
	// API routes
	apiGroup := router.Group("/api/v1")
	{
		// Client routes
		clients := apiGroup.Group("/clients")
		{
			clients.GET("", api.GetClients)
			clients.GET("/:id", api.GetClient)
			clients.POST("", api.CreateClient)
			clients.PUT("/:id", api.UpdateClient)
			clients.DELETE("/:id", api.DeleteClient)
			clients.GET("/:client_id/invoices", api.GetInvoicesByClient)
		}
		
		// Invoice routes
		invoices := apiGroup.Group("/invoices")
		{
			invoices.GET("", api.GetInvoices)
			invoices.GET("/:id", api.GetInvoice)
			invoices.POST("", api.CreateInvoice)
			invoices.PUT("/:id", api.UpdateInvoice)
			invoices.DELETE("/:id", api.DeleteInvoice)
		}
	}
	
	return router
}