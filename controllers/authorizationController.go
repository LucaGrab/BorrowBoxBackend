package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func isValueInArray(jsonArray string, value string, fieldName string) bool {
	// JSON-Array in ein Slice von Maps (JSON-Objekte) konvertieren
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonArray), &data); err != nil {
		panic(err)
	}

	// Durch das Array iterieren und überprüfen, ob der Wert im Feld vorhanden ist
	for _, item := range data {
		if fieldValue, ok := item[fieldName].(string); ok && fieldValue == value {
			return true
		}
	}
	return false
}

func Login(c *gin.Context) {
	var loginData map[string]string
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	email := loginData["email"]
	password := loginData["password"]
	client, err := database.NewMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	users := client.Database("borrowbox").Collection("users")

	filter := bson.M{
		"username": email,
		"password": password,
	}

	fmt.Println(password)

	var user models.User
	if err := users.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		print("invalid user")
		return
	}
	loginToken := user.ID
	c.JSON(http.StatusOK, gin.H{"loginToken": loginToken})
}
