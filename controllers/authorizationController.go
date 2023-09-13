package controllers

import (
	"BorrowBox/database"
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

	var loginData map[string]string // Erstelle eine Map, um E-Mail und Passwort zu speichern
	// Versuche, das JSON aus dem Request-Body in die loginData-Map zu binden
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Hier kannst du auf die Werte von E-Mail und Passwort zugreifen
	email := loginData["email"]
	password := loginData["password"]
	print(email)
	print(password)
	users, err := database.GetAllDcoumentsByCollection("users")
	if err != nil {
		return
	}
	// JSON-Array erstellen
	var jsonArray []string
	for _, bsonData := range users {
		jsonBytes, err := bson.MarshalExtJSON(bsonData, false, false)
		if err != nil {
			fmt.Println("Fehler bei der Umwandlung in JSON:", err)
			return
		}
		jsonString := string(jsonBytes)
		jsonArray = append(jsonArray, jsonString)
	}

	//finalJSON := "[" + strings.Join(jsonArray, ",") + "]"

	//fmt.Println(isValueInArray(finalJSON, email, "user"))

	//fmt.Println(finalJSON)

	loginToken := "99123455646372810987"

	c.JSON(http.StatusOK, gin.H{"loginToken": loginToken})

}
