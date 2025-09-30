package handlers

import (
	"net/http"

	"github.com/adem299/gow-commerce.git/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetProfile(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       currentUser.ID,
		"username": currentUser.Username,
	})
}
