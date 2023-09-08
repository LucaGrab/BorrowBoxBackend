package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Role  string `json:"role"`
	Email string `json:"email"`
}

type Rental struct {
	User   string    `json:"user"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Item   string    `json:"item"`
	Active bool      `json:"active"`
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	err := DeleteDocument("users", id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "User deleted"})
}

func getTagsById(itemId string) []string {
	tags := []string{}
	tagIds, err := getTagIdsByItemId(itemId)
	if err != nil {
		return nil
	}
	for _, tagIdMapping := range tagIds {
		tagId := tagIdMapping["tagId"]
		if oid, ok := tagId.(primitive.ObjectID); ok {
			tagIdString := oid.Hex()
			tag, err := getDocumentByID("tags", tagIdString)
			if err != nil {
				return nil
			}
			if tagName, ok := tag["name"].(string); ok {
				tags = append(tags, tagName)
			}
		}

	}
	return tags
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

func getDocumentByIDROute(c *gin.Context) {
	collection := c.Param("collection")
	id := c.Param("id")
	document, err := getDocumentByID(collection, id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	if collection == "items" {
		tags := getTagsById(id)
		document["tags"] = tags
	}
	c.IndentedJSON(200, document)
}

func getDocuments(c *gin.Context) {
	collection := c.Param("collection")
	documents, err := getAllDcoumentsByCollection(collection)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, documents)
}

func insertUser(c *gin.Context) {
	var newUser User // Ersetze YourDataStruct mit der tatsächlichen Struktur deiner Daten

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err := InsertDocument("users", newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "User inserted successfully"})
}

func insertRental(c *gin.Context) {
	var newRental Rental // Ersetze YourDataStruct mit der tatsächlichen Struktur deiner Daten

	if err := c.ShouldBindJSON(&newRental); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	newRental.Start = time.Now()
	newRental.Active = true

	err := InsertDocument("rentals", newRental)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert rental"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Rental inserted successfully"})
}

func updateUser(c *gin.Context) {
	id := c.Param("id") // ID des zu aktualisierenden Benutzers

	var updatedUser User // Ersetze YourDataStruct mit der tatsächlichen Struktur deiner Daten

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	updateData := bson.M{
		"$set": updatedUser, // Verwende die gesamte Struktur für das Update
	}

	err := UpdateDocument("users", id, updateData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func startGinServer() {

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8100"} // Add your frontend addresses here
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	r.Use(cors.New(config))

	r.GET("user/:id", userById)
	r.GET("documents/:collection", getDocuments)
	r.DELETE("user/:id", deleteUser)
	r.POST("user", insertUser)
	r.PUT("/user/:id", updateUser)
	r.GET("getDocumentByID/:collection/:id", getDocumentByIDROute)
	r.POST("startRental", insertRental)

	r.GET("/hello", func(c *gin.Context) { // bitte nicht löschen, ist gut zum testen
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.Run(":8080") // Starte den Gin-Server auf Port 8080
}
