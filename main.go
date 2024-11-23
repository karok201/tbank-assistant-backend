package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"tbank-assistant-backend/pkg/auth"
	"tbank-assistant-backend/pkg/books"
	"tbank-assistant-backend/pkg/common/db"
	"tbank-assistant-backend/pkg/user"
)

func main() {
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	r := gin.Default()
	h := db.Init(dbUrl)
	e := db.InitElastic()

	books.RegisterRoutes(r, h, e)
	auth.RegisterRoutes(r, h)
	user.RegisterRoutes(r, h, e)

	r.Run(port)
}
