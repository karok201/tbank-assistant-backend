package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tbank-assistant-backend/pkg/common/models"
	"tbank-assistant-backend/pkg/common/utils"
)

type GetUserResponseBody struct {
	ID       int    `json:"int"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

func (h handler) GetUser(c *gin.Context) {
	var userFound models.User

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Where("id=?", userId).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	response := GetUserResponseBody{
		ID:       userId,
		Email:    userFound.Email,
		Username: userFound.Username,
		Phone:    userFound.Phone,
	}

	c.JSON(http.StatusOK, response)
}
