package user

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tbank-assistant-backend/pkg/common/models"
	"tbank-assistant-backend/pkg/common/utils"
)

type GetOperationsRequestBody struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type GetOperationsResponseBody struct {
	Id       string `json:"id"`
	Type     int    `json:"type"`
	Category string `json:"category"`
	Date     string `json:"date"`
	Amount   int    `json:"amount"`
	Comment  string `json:"comment"`
}

func (h handler) GetOperationsPeriod(c *gin.Context) {
	var userFound models.User

	body := GetOperationsRequestBody{}

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

	query := map[string]interface{}{
		"size": 10000,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"terms": map[string]interface{}{
							"user_id": []string{strconv.Itoa(userId)},
						},
					},
					map[string]interface{}{
						"range": map[string]interface{}{
							"date": map[string]interface{}{
								"gte": body.StartDate,
								"lte": body.EndDate,
							},
						},
					},
				},
			},
		},
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

	var documents []map[string]interface{}
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		doc := map[string]interface{}{
			"id":       hit.(map[string]interface{})["_id"].(string),
			"type":     int(source.(map[string]interface{})["type"].(float64)),
			"amount":   int(source.(map[string]interface{})["amount"].(float64)),
			"date":     source.(map[string]interface{})["date"].(string),
			"category": int(source.(map[string]interface{})["category"].(float64)),
		}
		documents = append(documents, doc)
	}

	c.JSON(http.StatusOK, documents)
}
