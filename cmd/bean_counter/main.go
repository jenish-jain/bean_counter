package main

import (
	"bean_counter/internal/server"
	"log"

	"github.com/jenish-jain/logger"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	logger.Init("debug")

	server := server.NewServer()
	server.Run()

}
