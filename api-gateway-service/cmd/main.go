package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prateek041/api-gateway-service/pkg/middleware"
)

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware()) // Enable CORS for development

	router.GET("/health", healthHandler)

	// Public routes (no authentication required)
	public := router.Group("/auth")
	{
		public.Any("/login", proxyToService("http://user-management-service:8081"))
		public.Any("/users", proxyToService("http://user-management-service:8081")) // Assuming signup is here
	}

	// Protected routes (authentication required)
	protected := router.Group("/")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.Any("/products/*path", proxyToService("http://product-presentation-service:8080"))
		// Add more protected routes for other services here
	}

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = ":8000"
	} else {
		port = ":" + port
	}

	log.Printf("API Gateway Service listening on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Error starting gateway server: %v", err)
	}
}

func proxyToService(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		log.Printf("URL: %s", c.Request.URL)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
