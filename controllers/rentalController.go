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

	_, err := database.InsertDocument("rentals", newRental)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert rental"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Rental inserted successfully"})
}

func EndRental(c *gin.Context) {
	var returnData models.ReturnData

	if err := c.ShouldBindJSON(&returnData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	itemId, err := primitive.ObjectIDFromHex(returnData.ItemId.Hex())
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	rental := database.GetActiveRentalByItemId(itemId)
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

	//update location
	println("updating location...")
	itemData, err := database.GetDocumentByID("items", itemId.Hex())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location (getting item)"})
		return
	}
	println(itemData)

	var item models.ItemForInsert

	// Convert bson.M to raw BSON data
	itemBSONData, err := bson.Marshal(itemData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert item data to BSON"})
		return
	}

	// Unmarshal raw BSON data into the item struct
	if err := bson.Unmarshal(itemBSONData, &item); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal item data"})
		return
	}

	item.Location = returnData.Location
	//print item data
	println(item.ID.Hex())
	println(item.Location)
	println(item.Name)
	println(item.Description)

	updateData = bson.M{
		"$set": item, // Verwende die gesamte Struktur für das Update
	}
	err = database.UpdateDocument("items", itemId.Hex(), updateData)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location (updating item)"})
		return
	}

	println("...location updated")

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Rental updated successfully"})
}
