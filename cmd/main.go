package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"tbank-assistant-backend/pkg/auth"
	"tbank-assistant-backend/pkg/books"
	"tbank-assistant-backend/pkg/common/db"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	port := viper.Get("PORT").(string)
	dbUrl := viper.Get("DB_URL").(string)

	r := gin.Default()
	h := db.Init(dbUrl)

	books.RegisterRoutes(r, h)
	auth.RegisterRoutes(r, h)

	r.Run(port)
}
