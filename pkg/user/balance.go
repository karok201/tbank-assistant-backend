package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tbank-assistant-backend/pkg/common/models"
	"tbank-assistant-backend/pkg/common/utils"
)

type GetUserBalanceResponse struct {
	Amount string `json:"amount"`
}

func (h handler) GetUserBalance(c *gin.Context) {
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

	sumSpent, _ := getSumByType(userId, 0, h)
	sumGot, _ := getSumByType(userId, 1, h)

	response := GetUserBalanceResponse{
		Amount: fmt.Sprintf("%d", sumGot-sumSpent),
	}

	c.JSON(http.StatusOK, response)
	return
}

func getSumByType(userId int, transactionType int, h handler) (int, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"terms": map[string]interface{}{
							"user_id": []string{strconv.Itoa(userId)},
						},
					},
					map[string]interface{}{
						"terms": map[string]interface{}{
							"type": []string{strconv.Itoa(transactionType)},
						},
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"aggregations": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "user_id",
					"size":  10000,
				},
				"aggs": map[string]interface{}{
					"spent": map[string]interface{}{
						"sum": map[string]interface{}{
							"field": "amount",
						},
					},
				},
			},
		},
		"size": 0,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := h.ES.Search(
		h.ES.Search.WithIndex("user_operations"),
		h.ES.Search.WithBody(&buf),
		h.ES.Search.WithTrackTotalHits(true), // Track the total number of hits
		h.ES.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	result := r["aggregations"].(map[string]interface{})["aggregations"].(map[string]interface{})["buckets"].([]interface{})[0].(map[string]interface{})["spent"].(map[string]interface{})["value"].(float64)

	return int(result), err
}
