package controllers

import (
    "github.com/gin-gonic/gin"
    "mobile-backend-go/models"
    "mobile-backend-go/database"
    "net/http"
    "log"
    "strconv"
)

// GetClient returns a client by ID
// @Summary Get a client by ID
// @Description Fetch a client by its ID
// @Tags Clients
// @Security BearerAuth
// @Produce  json
// @Param id path int true "Client ID"
// @Success 200 {object} models.Client
// @Failure 400 {object} map[string]string "Invalid client ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Client not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/clients/{id} [get]
func GetClient(c *gin.Context) {
    userID := c.MustGet("userID").(uint)
    clientID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
        return
    }

    var client models.Client
    if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
        return
    }

    c.JSON(http.StatusOK, client)
}

// GetClients returns list of all clients
// @Summary Get list of clients
// @Description Get all clients for the authenticated user
// @Tags Clients
// @Security BearerAuth
// @Produce  json
// @Success 200 {array} models.Client
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/clients [get]
func GetClients(c *gin.Context) {
    userID := c.MustGet("userID").(uint)

    var clients []models.Client
    if err := database.DB.Where("user_id = ?", userID).Find(&clients).Error; err != nil {
        log.Printf("Failed to fetch clients for userID %v: %v", userID, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clients"})
        return
    }

    c.JSON(http.StatusOK, clients)
}

// AddClient adds a new client
// @Summary Add a new client
// @Description Create a new client for the authenticated user
// @Tags Clients
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param client body models.Client true "Client data"
// @Success 201 {object} models.Client
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/clients [post]
func AddClient(c *gin.Context) {
    userID := c.MustGet("userID").(uint)

    var newClient models.Client
    if err := c.ShouldBindJSON(&newClient); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newClient.UserID = userID
    if err := database.DB.Create(&newClient).Error; err != nil {
        log.Printf("Failed to add client: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add client"})
        return
    }

    c.JSON(http.StatusCreated, newClient)
}

// UpdateClient updates client data
// @Summary Update a client
// @Description Update a client's details
// @Tags Clients
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Client ID"
// @Param client body models.Client true "Client data"
// @Success 200 {object} map[string]string "Client updated successfully"
// @Failure 400 {object} map[string]string "Invalid client ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Client not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/clients/{id} [put]
func UpdateClient(c *gin.Context) {
    userID := c.MustGet("userID").(uint)
    clientID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
        return
    }

    var client models.Client
    if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
        return
    }

    if err := c.ShouldBindJSON(&client); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := database.DB.Save(&client).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Client updated successfully"})
}

// DeleteClient deletes a client
// @Summary Delete a client
// @Description Delete a client by its ID
// @Tags Clients
// @Security BearerAuth
// @Param id path int true "Client ID"
// @Success 200 {object} map[string]string "Client deleted successfully"
// @Failure 400 {object} map[string]string "Invalid client ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Client not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/clients/{id} [delete]
func DeleteClient(c *gin.Context) {
    userID := c.MustGet("userID").(uint)
    clientID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
        return
    }

    var client models.Client
    if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
        return
    }

    if err := database.DB.Delete(&client).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}