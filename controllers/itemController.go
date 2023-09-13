package controllers

import (
	"BorrowBox/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserItems(c *gin.Context) {
	id := c.Param("id")
	rentals, err := database.GetDocumentsByCollectionFiltered("rentals", "userId", id, true)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	itemIds := []string{}
	for _, rental := range rentals {
		if rental["active"] == true {
			itemIds = append(itemIds, rental["itemId"].(primitive.ObjectID).Hex())
		}
	}
	items := []bson.M{}
	for _, itemId := range itemIds {
		item, err := database.GetDocumentByID("items", itemId)
		if err != nil {
			c.IndentedJSON(404, gin.H{"message": err.Error()})
			return
		}
		items = append(items, item)
	}
	c.IndentedJSON(200, items)
}
