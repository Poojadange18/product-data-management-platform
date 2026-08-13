package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"product-data-management-platform/config"
	"product-data-management-platform/handler"
	"product-data-management-platform/middleware"
	"product-data-management-platform/repository"
	"product-data-management-platform/service"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	products := handler.NewProductHandler(service.NewProductService(repository.NewProductRepository(db)))
	authService := service.NewAuthService(repository.NewUserRepository(db), env("JWT_SECRET", "change-me-in-production"))
	auth := handler.NewAuthHandler(authService)
	r := gin.Default()
	origin := env("CORS_ALLOWED_ORIGIN", "http://localhost:8080")
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.Static("/assets", "./frontend")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/login", "./frontend/login.html")
	r.StaticFile("/register", "./frontend/register.html")
	r.StaticFile("/openapi.yaml", "./docs/openapi.yaml")
	r.StaticFile("/docs", "./docs/index.html")
	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)
	protected := r.Group("/products")
	protected.Use(middleware.RequireAuth(authService))
	protected.POST("", products.CreateProduct)
	protected.GET("", products.GetProducts)
	protected.GET("/:id", products.GetProductByID)
	protected.PUT("/:id", products.UpdateProduct)
	protected.DELETE("/:id", products.DeleteProduct)
	if err := r.Run(":" + env("SERVER_PORT", "8080")); err != nil {
		log.Fatal(err)
	}
}
