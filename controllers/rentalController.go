package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func InsertRental(c *gin.Context) {
	var newRental models.Rental // Ersetze YourDataStruct mit der tatsächlichen Struktur deiner Daten

	if err := c.ShouldBindJSON(&newRental); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	newRental.Start = time.Now()
	newRental.Active = true

	err := database.InsertDocument("rentals", newRental)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert rental"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Rental inserted successfully"})
}

func EndRental(c *gin.Context) {
	itemId := c.Param("itemId")
	objectId, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	rental := database.GetActiveRentalByItemId(objectId)

	rental.End = time.Now()
	rental.Active = false

	updateData := bson.M{
		"$set": rental, // Verwende die gesamte Struktur für das Update
	}

	err = database.UpdateDocument("rentals", rental.Id.Hex(), updateData)
	if err != nil {
		//print error
		println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rental"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Rental updated successfully"})
}
