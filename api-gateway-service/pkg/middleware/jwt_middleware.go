package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	publicKey     []byte
	publicKeyOnce sync.Once
	publicKeyErr  error
)

// TODO: This should be configured through Environment variables.
const publicKeyURL = "http://user-management-service:8081/auth/public-key" // Kubernetes Service name

func getPublicKey() ([]byte, error) {
	publicKeyOnce.Do(func() {
		resp, err := http.Get(publicKeyURL)
		if err != nil {
			publicKeyErr = fmt.Errorf("failed to fetch public key: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			publicKeyErr = fmt.Errorf("failed to fetch public key, status code: %d", resp.StatusCode)
			return
		}

		publicKey, err = io.ReadAll(resp.Body)
		if err != nil {
			publicKeyErr = fmt.Errorf("failed to read public key: %w", err)
			return
		}
		log.Println("Public key fetched successfully")
	})
	return publicKey, publicKeyErr
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := tokenParts[1]

		key, err := getPublicKey()
		if err != nil {
			log.Printf("Error getting public key: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
			return rsaPublicKey, err
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			// You can extract user information from claims if needed
			c.Set("claims", claims)
			c.Next() // Proceed to the next handler
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

