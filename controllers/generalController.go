package controllers

import (
	"BorrowBox/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

func GetDocumentByIDROute2(c *gin.Context) {
	collection := c.Param("collection")
	id := c.Param("id")
	formattedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": formattedId},
		},
		{
			"$lookup": bson.M{
				"from":         "itemTag",
				"localField":   "_id",
				"foreignField": "itemId",
				"as":           "itemTags",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "rentals",
				"localField":   "_id",
				"foreignField": "itemId",
				"as":           "rentals",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "tags",
				"localField":   "itemTags.tagId",
				"foreignField": "_id",
				"as":           "tags",
			},
		},
		{
			"$unwind": "$tags", // Entfalte das "tags"-Array
		},
		{
			"$addFields": bson.M{
				"rentals": bson.M{
					"$filter": bson.M{
						"input": "$rentals",
						"as":    "rental",
						"cond":  bson.M{"$eq": []interface{}{"$$rental.active", true}},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":         "$_id",
				"description": bson.M{"$first": "$description"},
				"location":    bson.M{"$first": "$location"},
				"name":        bson.M{"$first": "$name"},
				"tagNames":    bson.M{"$push": "$tags.name"}, // Extrahiere die Tag-Namen in ein Array
				"rentals":     bson.M{"$first": "$rentals"},  // Behalte das gefilterte "rentals"-Array bei
			},
		},
		{
			"$project": bson.M{
				"_id":         0, // Ausblenden der _id-Felder
				"description": 1,
				"location":    1,
				"name":        1,
				"tagNames":    1, // Das Array mit Tag-Namen beibehalten
				"rentals":     1, // Das gefilterte "rentals"-Array beibehalten
			},
		},
	}

	documents, err := database.NewDBAggregation(collection, pipeline)

	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, documents)
}
