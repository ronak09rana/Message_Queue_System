package main

import (
	"fmt"
	"log"
	"message_queue_system/bootstrap"
	"message_queue_system/db"
	"message_queue_system/domain"
	"message_queue_system/rabbitmq"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error: %v, unable to load .env file", err)
		return
	}

	dbCred := getCredentials()
	err = db.Init(dbCred)
	if err != nil {
		log.Printf("Error: %v, unable to connect to database", err.Error())
		return
	}

	err = rabbitmq.Connect()
	if err != nil {
		log.Printf("Error: %v, unable to init rabbitmq", err.Error())
		return
	}

	r := gin.Default()
	apiGroup := r.Group("/api")
	bootstrap.Init(apiGroup)

	PORT := os.Getenv("PORT")
	err = r.Run(fmt.Sprintf(":%v", PORT))
	if err != nil {
		log.Printf("Error: %v, unable to run server", err.Error())
	}
	log.Printf("Server Running on PORT: %v", PORT)
}

func getCredentials() domain.DBCred {
	return domain.DBCred{
		Username: os.Getenv("username"),
		Password: os.Getenv("password"),
		Hostname: os.Getenv("hostname"),
		DBName:   os.Getenv("dbname"),
	}
}
