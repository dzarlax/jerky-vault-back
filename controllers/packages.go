package controllers

import (
    "github.com/gin-gonic/gin"
    "mobile-backend-go/models"
    "mobile-backend-go/database"
    "net/http"
)

// GetPackages возвращает список всех упаковок
// @Summary Get list of packages
// @Description Get all packages for the authenticated user
// @Tags Packages
// @Security BearerAuth
// @Produce  json
// @Success 200 {array} models.Package
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/packages [get]
func GetPackages(c *gin.Context) {
    userID := c.MustGet("userID").(uint)

    var packages []models.Package
    if err := database.DB.Where("user_id = ?", userID).Find(&packages).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch packages"})
        return
    }

    c.JSON(http.StatusOK, packages)
}

// AddPackage добавляет новую упаковку
// @Summary Add a new package
// @Description Create a new package for the authenticated user
// @Tags Packages
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param package body models.Package true "Package data"
// @Success 201 {object} models.Package
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/packages [post]
func AddPackage(c *gin.Context) {
    userID := c.MustGet("userID").(uint)

    var newPackage models.Package
    if err := c.ShouldBindJSON(&newPackage); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newPackage.UserID = userID
    if err := database.DB.Create(&newPackage).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add package"})
        return
    }

    c.JSON(http.StatusCreated, newPackage)
}