package user

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"tbank-assistant-backend/pkg/common/models"
	"tbank-assistant-backend/pkg/common/utils"
)

type CreateOperationsChequeByRequestBody struct {
	ChequeData string `json:"cheque_data"`
}

type CreateOperationsChequeByResponseBody struct {
	Result string `json:"result"`
}

func (h handler) CreateOperationsByCheque(c *gin.Context) {
	var userFound models.User

	var requestBody CreateOperationsChequeByRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
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

	url := "https://proverkacheka.com/api/v1/check/get"

	data := map[string]interface{}{
		"token": "30176.dwhLbhcjry8B1FrRy",
		"qrraw": requestBody.ChequeData,
	}

	// Сериализация структуры в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	// Создание нового запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %s", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %s", err)
	}

	// Парсинг JSON из ответа
	var responseJSON map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseJSON); err != nil {
		log.Fatalf("Error unmarshaling JSON: %s", err)
	}

	// Извлечение объекта "data"
	data, ok := responseJSON["data"].(map[string]interface{})["json"].(map[string]interface{})
	if !ok {
		log.Fatalf("Error: 'data' is not an object")
	}

	totalSum := int(data["totalSum"].(float64)) / 100
	date := data["dateTime"].(string)[:10]
	category := 0
	operationType := int(data["operationType"].(float64))
	comment := ""

	document := struct {
		UserId   int    `json:"user_id"`
		Type     int    `json:"type"`
		Category int    `json:"category"`
		Date     string `json:"date"`
		Amount   int    `json:"amount"`
		Comment  string `json:"comment"`
	}{
		UserId:   userId,
		Type:     operationType,
		Category: category,
		Date:     date,
		Amount:   totalSum,
		Comment:  comment,
	}

	jsonDocument, err := json.Marshal(document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal document"})
		return
	}

	res, err := h.ES.Index("user_operations", bytes.NewReader(jsonDocument))

	println(res.String())
	if err != nil {
		c.JSON(res.StatusCode, gin.H{"error": res.String()})
	}

	c.JSON(http.StatusOK, gin.H{"res": "success"})
}
