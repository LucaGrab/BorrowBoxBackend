package main

import (
	"BorrowBox/database"
	"BorrowBox/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvVariables()
	database.Connect()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8100"} // Add your frontend addresses here
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	r.Use(cors.New(config))

	routes.Setup(r)

	r.Run(":8080") // Starte den Gin-Server auf Port 8080

}

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
