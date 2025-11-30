package main

import (
	"github.com/adem299/gow-commerce.git/database"
	"github.com/adem299/gow-commerce.git/handlers"
	"github.com/adem299/gow-commerce.git/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := database.ConnectDatabase()

	h := handlers.NewHandler(db)

	router := gin.Default()

	// cors configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Public routes

	public := router.Group("/auth")
	{
		public.POST("/register", h.Register)
		public.POST(("/login"), h.Login)
	}

	router.GET("/products", h.GetProducts)
	router.GET("/products/:id", h.GetProductByID)

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware(db))
	{
		api.GET("/profile", h.GetProfile)

		// Cart routes
		api.POST("/cart", h.AddToCart)
		api.GET("/cart", h.GetCart)

		// Order routes
		api.POST("/orders", h.CreateOrder)
		api.GET("/orders", h.GetOrders)

		// Admin routes
		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.AdminMiddleware())
		{
			adminRoutes.POST("/products", h.CreateProduct)
			adminRoutes.PUT("/products/:id", h.UpdateProduct)
			adminRoutes.DELETE("/products/:id", h.DeleteProduct)
			adminRoutes.GET("/orders", h.GetAllOrders)
			adminRoutes.PUT("/orders/:id/status", h.UpdateOrderStatus)
		}
	}

	router.Run("localhost:8080")
}
