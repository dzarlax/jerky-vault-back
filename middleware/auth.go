package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Claims struct {
	UserID uint `json:"userID"`
	jwt.StandardClaims
}

// validateJWTSecret checks that JWT_SECRET is set and has sufficient length
func validateJWTSecret() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is not set")
	}
	if len(secret) < 16 {
		return fmt.Errorf("JWT_SECRET must be at least 16 characters long for security")
	}
	return nil
}

// init performs validation when loading the package
func init() {
	if err := validateJWTSecret(); err != nil {
		log.Fatalf("JWT configuration error: %v", err)
	}
}

// JWTMiddleware validates JWT token in requests and sets userID in context
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix from token
		if len(tokenString) > 7 && strings.ToUpper(tokenString[:7]) == "BEARER " {
			tokenString = tokenString[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Check signing method to protect against "none algorithm" attack
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			// Check specifically for token expiration
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
					c.Abort()
					return
				}
			}
			log.Printf("Invalid token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Additional check for token expiration
		if claims.ExpiresAt > 0 && time.Now().Unix() > claims.ExpiresAt {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		// Set userID in context for use in controllers
		c.Set("userID", claims.UserID)
		c.Next()
	}
}