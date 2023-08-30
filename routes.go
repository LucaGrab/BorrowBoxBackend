package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id    string `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

func userById(c *gin.Context) {
	id := c.Param("id")
	user, err := getDocumentByID("users", id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, user)
}

func getUsers(c *gin.Context) {
	users, err := getAllDcoumentsByCollection("users")
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, users)
}

func startGinServer() {

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8100"} // Add your frontend addresses here
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	r.Use(cors.New(config))

	r.GET("user/:id", userById)
	r.GET("users", getUsers)

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.Run(":8080") // Starte den Gin-Server auf Port 8080
}
