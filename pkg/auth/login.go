package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"tbank-assistant-backend/pkg/common/models"
	"tbank-assistant-backend/pkg/common/utils"
)

func (h handler) LoginUser(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user := models.User{}

	result := h.DB.Where("email = ?", credentials.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User is not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Println(user.CheckPassword(credentials.Password))
	if user.CheckPassword(credentials.Password) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
