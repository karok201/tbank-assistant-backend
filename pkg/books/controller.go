package books

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
	ES *elasticsearch.Client
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB, es *elasticsearch.Client) {
	h := &handler{
		DB: db,
		ES: es,
	}

	routes := r.Group("/books")
	routes.POST("/", h.AddBook)
	routes.GET("/", h.GetBooks)
	routes.GET("/:id", h.GetBook)
	routes.PUT("/:id", h.UpdateBook)
	routes.DELETE("/:id", h.DeleteBook)
}
