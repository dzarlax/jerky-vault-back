package controllers

import (
    "github.com/gin-gonic/gin"
    "mobile-backend-go/database"
    "mobile-backend-go/models"
    "net/http"
)

// CreateCookingSession creates a new cooking session
// @Summary Create a new cooking session
// @Description Create a new cooking session with details
// @Tags Cooking Sessions
// @Accept  json
// @Produce  json
// @Param session body models.CookingSessionCreateDTO true "Cooking Session data"
// @Success 201 {object} models.CookingSession
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/cooking_sessions [post]
func CreateCookingSession(c *gin.Context) {
    userID := c.MustGet("userID").(uint)

    var requestData models.CookingSessionCreateDTO
    if err := c.ShouldBindJSON(&requestData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Create CookingSession model from DTO
    newSession := models.CookingSession{
        RecipeID: requestData.RecipeID,
        Date:     requestData.Date,
        Yield:    requestData.Yield,
        UserID:   userID,
    }

    if err := database.DB.Create(&newSession).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cooking session"})
        return
    }

    // Preload Recipe for response
    if err := database.DB.Preload("Recipe").First(&newSession, newSession.ID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load cooking session"})
        return
    }

    c.JSON(http.StatusCreated, newSession)
}

// GetCookingSessions returns list of all cooking sessions
// @Summary Get list of cooking sessions
// @Description Get all cooking sessions available for the authenticated user
// @Tags Cooking Sessions
// @Security BearerAuth
// @Produce  json
// @Success 200 {array} models.CookingSession
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/cooking_sessions [get]
func GetCookingSessions(c *gin.Context) {
    userID := c.MustGet("userID").(uint)
    var sessions []models.CookingSession

    if err := database.DB.Where("user_id = ?", userID).Preload("Ingredients.Ingredient").Find(&sessions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cooking sessions"})
        return
    }

    c.JSON(http.StatusOK, sessions)
}