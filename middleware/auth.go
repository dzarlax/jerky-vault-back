package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    "net/http"
    "os"
    "log"
)

type Claims struct {
    UserID uint `json:"userID"`
    jwt.StandardClaims
}

// JWTMiddleware проверяет JWT токен в запросах и устанавливает userID в контекст
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            log.Println("Authorization header missing")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
            c.Abort()
            return
        }

        // Убираем префикс "Bearer " из токена
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            log.Printf("Invalid token: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Логируем извлеченный userID
        log.Printf("Extracted userID from token: %v", claims.UserID)

        // Установка userID в контекст для использования в контроллерах
        c.Set("userID", claims.UserID)
        c.Next()
    }
}