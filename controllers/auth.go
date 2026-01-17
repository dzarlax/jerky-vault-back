package controllers

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    "mobile-backend-go/models"
    "mobile-backend-go/database"
    "mobile-backend-go/utils"
    "net/http"
    "os"
    "time"
)

// Definition of Claims structure for JWT tokens
type Claims struct {
    UserID uint `json:"userID"` // Changed to UserID
    jwt.StandardClaims
}

// Definition of structures for registration and login requests
type RegisterPayload struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginPayload struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// Register creates a new user
// @Summary Register a new user
// @Description Create a new user with a username and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body RegisterPayload true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
    var payload RegisterPayload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := utils.HashPassword(payload.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    user := models.User{Username: payload.Username, Password: hashedPassword}
    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
        return
    }

    c.JSON(http.StatusCreated, user)
}

// Login authenticates a user and returns a JWT
// @Summary Login a user
// @Description Authenticate a user and return a JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body LoginPayload true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
    var payload LoginPayload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := database.DB.Where("username = ?", payload.Username).First(&user).Error; err != nil {
        // Do not reveal information about user existence
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if !utils.CheckPassword(user.Password, payload.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Token is valid for 24 hours (instead of 72)
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
            IssuedAt:  time.Now().Unix(),
            Subject:   fmt.Sprintf("user:%d", user.ID),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": tokenString,
        "expires_at": expirationTime,
    })
}