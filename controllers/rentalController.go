package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"github.com/gin-gonic/gin"
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
