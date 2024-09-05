package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
    "mobile-backend-go/models"
    "mobile-backend-go/database"
    "mobile-backend-go/utils"
    "net/http"
    "os"
    "time"
)

// Определение структуры Claims для JWT токенов
type Claims struct {
    UserID uint `json:"userID"` // Изменено на UserID
    jwt.StandardClaims
}

// Определение структур для запросов регистрации и входа
type RegisterPayload struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginPayload struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// Register создает нового пользователя
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

// Login аутентифицирует пользователя и возвращает JWT
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
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    if !utils.CheckPassword(user.Password, payload.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        return
    }

    expirationTime := time.Now().Add(72 * time.Hour)
    claims := &Claims{
        UserID: user.ID, // Теперь используем UserID
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}