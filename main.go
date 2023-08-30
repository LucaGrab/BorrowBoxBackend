package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvVariables()
	startGinServer()
}

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
