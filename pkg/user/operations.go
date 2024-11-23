package user

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"tbank-assistant-backend/pkg/common/models"
	"tbank-assistant-backend/pkg/common/utils"
)

type CreateOperationRequestBody struct {
	Type     int    `json:"type"`
	Category string `json:"category"`
	Date     string `json:"date"`
	Amount   int    `json:"amount"`
	Comment  string `json:"comment"`
}

func (h handler) CreateOperation(c *gin.Context) {
	var userFound models.User

	body := CreateOperationRequestBody{}

	if err := c.BindJSON(&body); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Where("id=?", userId).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	document := struct {
		UserId   int    `json:"user_id"`
		Type     int    `json:"type"`
		Category string `json:"category"`
		Date     string `json:"date"`
		Amount   int    `json:"amount"`
		Comment  string `json:"comment"`
	}{
		UserId:   userId,
		Type:     body.Type,
		Category: body.Category,
		Date:     body.Date,
		Amount:   body.Amount,
		Comment:  body.Comment,
	}

	jsonDocument, err := json.Marshal(document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal document"})
		return
	}

	res, err := h.ES.Index("user_operations", bytes.NewReader(jsonDocument))
	if err != nil {
		c.JSON(res.StatusCode, gin.H{"error": res.String()})
	}

	c.JSON(http.StatusOK, gin.H{"res": "success"})
}
