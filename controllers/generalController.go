package controllers

import (
	"BorrowBox/database"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDocuments(c *gin.Context) {
	collection := c.Param("collection")
	documents, err := database.GetAllDcoumentsByCollection(collection)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	if collection == "items" {
		for i := range documents {
			document := documents[i]
			// Zugriff auf das "_id" Feld des Dokuments
			id := document["_id"].(primitive.ObjectID) // Annahme: Verwendung von BSON für MongoDB

			idString := id.Hex()
			rentals, err := database.GetDocumentsByCollectionFiltered("rentals", "itemId", idString, true)
			if err != nil {
				c.IndentedJSON(404, gin.H{"message": err.Error()})
				return
			}
			document["available"] = true
			for _, rental := range rentals {
				if rental["active"] == true {
					document["available"] = false
					break
				}
			}
			tags := GetTagsById(idString)
			document["tags"] = tags
		}
	}

	c.IndentedJSON(200, documents)
}

func GetDocumentByIDROute(c *gin.Context) {
	collection := c.Param("collection")
	id := c.Param("id")
	document, err := database.GetDocumentByID(collection, id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	if collection == "items" {
		tags := GetTagsById(id)
		document["tags"] = tags
		currentRentals, err := database.GetDocumentsByCollectionFiltered2("rentals", "itemId", id, true, "active", true, false)
		if err != nil {
			c.IndentedJSON(404, gin.H{"message": err.Error()})
			return
		}
		if len(currentRentals) == 0 {
			document["available"] = true
		} else {
			document["available"] = false
			currentRenterId := currentRentals[0]["userId"].(primitive.ObjectID).Hex()
			currentRenter, err := database.GetDocumentByID("users", currentRenterId)
			if err != nil {
				c.IndentedJSON(404, gin.H{"message": err.Error()})
				return
			}
			document["currentRenter"] = currentRenter["email"]
			fmt.Println("Das Array ist nicht leer.")
		}
	}
	c.IndentedJSON(200, document)
}