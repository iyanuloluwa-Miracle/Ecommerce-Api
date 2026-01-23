package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/iyanuloluwa-Miracle/ecommerce-api/internal/controller"
	"github.com/iyanuloluwa-Miracle/ecommerce-api/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Health check endpoint
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "API is running"})
	})

	// Authentication routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Product routes (public)
	products := api.Group("/products")
	{
		products.GET("", controllers.GetProducts)
		products.GET("/:id", controllers.GetProduct)
	}

	// Cart routes (authenticated)
	cart := api.Group("/cart", middleware.JWTAuth())
	{
		cart.POST("/add", controllers.AddToCart)
		cart.DELETE("/remove", controllers.RemoveFromCart)
		cart.GET("", controllers.ViewCart)
	}

	// Checkout routes (authenticated)
	checkout := api.Group("/checkout", middleware.JWTAuth())
	{
		checkout.POST("", controllers.Checkout)
	}

	// Order routes (authenticated)
	orders := api.Group("/orders", middleware.JWTAuth())
	{
		orders.GET("", controllers.GetOrders)
	}

	// Stripe webhook (no auth required, uses Stripe signature verification)
	api.POST("/webhook/stripe", controllers.WebhookHandler)

	// Admin routes (authenticated + admin role)
	admin := api.Group("/admin", middleware.JWTAuth(), middleware.AdminOnly())
	{
		admin.POST("/products", controllers.CreateProduct)
		admin.PUT("/products/:id", controllers.UpdateProduct)
		admin.DELETE("/products/:id", controllers.DeleteProduct)
	}
}
