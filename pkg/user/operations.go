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
	Type     string `json:"type"`
	Category string `json:"category"`
	Date     string `json:"date"`
	Amount   string `json:"amount"`
	Comment  string `json:"comment"`
}

func (h handler) CreateOperation(c *gin.Context) {
	var userFound models.User

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
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Price int    `json:"price"`
	}{
		Id:    1,
		Name:  "Foo",
		Price: 10,
	}

	jsonDocument, err := json.Marshal(document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal document"})
		return
	}

	res, err := h.ES.Index("index_name", bytes.NewReader(jsonDocument))

	c.JSON(http.StatusOK, res)
}
