package db

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func InitElastic() *elasticsearch.Client {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		CloudID: "My_deployment:ZXVyb3BlLXdlc3Q0LmdjcC5lbGFzdGljLWNsb3VkLmNvbTo0NDMkYjJiZDEyNGRhNTIzNGVlYmFhMmNjZGUwNmUyMmQ5MjckODIyZTVhM2U1Y2QxNGQ5M2FkM2UzYTEyNzVlODVhOGU=",
		APIKey:  "bHdrWVY1TUJra3BBRmVMTWdEMUs6T2o3d0l2ZHlUb2VDYTFaOXlWWXpBdw==",
	})
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
