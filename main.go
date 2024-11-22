package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"tbank-assistant-backend/pkg/auth"
	"tbank-assistant-backend/pkg/books"
	"tbank-assistant-backend/pkg/common/db"
)

func main() {
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	r := gin.Default()
	h := db.Init(dbUrl)

	books.RegisterRoutes(r, h)
	auth.RegisterRoutes(r, h)

	r.Run(port)
}
