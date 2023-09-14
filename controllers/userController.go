package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsers(c *gin.Context) {
	collection := "users"
	documents, err := database.GetAllDcoumentsByCollection(collection)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, documents)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := database.DeleteDocument("users", id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "User deleted"})
}

func UserById(c *gin.Context) {
	id := c.Param("id")
	user, err := database.GetDocumentByID("users", id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, user)
}

func InsertUser(c *gin.Context) {
	var newUser models.User // Ersetze YourDataStruct mit der tatsächlichen Struktur deiner Daten

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err := database.InsertDocument("users", newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "User inserted successfully"})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id") // ID des zu aktualisierenden Benutzers

	var updatedUser models.User // Ersetze YourDataStruct mit der tatsächlichen Struktur deiner Daten

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	updateData := bson.M{
		"$set": updatedUser, // Verwende die gesamte Struktur für das Update
	}

	err := database.UpdateDocument("users", id, updateData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
