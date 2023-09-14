package controllers

import (
	"BorrowBox/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetActiveUserItems(c *gin.Context) {
	id := c.Param("id")
	formattedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	pipeline := []bson.M{
		{
			"$match": bson.M{"userId": formattedId, "active": true},
		},

		{
			"$lookup": bson.M{
				"from":         "items",
				"localField":   "itemId",
				"foreignField": "_id",
				"as":           "items",
			},
		},
		{
			"$project": bson.M{
				"_id":   0,
				"items": "$items",
			},
		},
	}
	documents, err := database.NewDBAggregation("rentals", pipeline)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	document := documents[0] //warum hat documents zusätzlich ein array außen in der antwort - so ist es trotzdem ein array von items

	c.IndentedJSON(200, document)
}

/*
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
*/
func GetItemByIdWithAllRentals(c *gin.Context) {
	collection := "items"
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
			"$unwind": "$rentals", // Entfalte das "rentals"-Array
		},
		{
			"$group": bson.M{
				"_id":         "$_id",
				"description": bson.M{"$first": "$description"},
				"location":    bson.M{"$first": "$location"},
				"name":        bson.M{"$first": "$name"},
				"tagNames":    bson.M{"$push": "$tags.name"},   // Extrahiere die Tag-Namen in ein Array
				"rentals":     bson.M{"$addToSet": "$rentals"}, // Sammle alle Rentals und entferne Duplikate
			},
		},
		{
			"$project": bson.M{
				"_id":         0, // Ausblenden der _id-Felder
				"description": 1,
				"location":    1,
				"name":        1,
				"tagNames":    1, // Das Array mit Tag-Namen beibehalten
				"rentals":     1, // Die "rentals" beibehalten
			},
		},
	}

	documents, err := database.NewDBAggregation(collection, pipeline)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	document := documents[0]
	c.IndentedJSON(200, document)
}

func GetItemByIdWithTheActiveRental(c *gin.Context) {
	collection := "items"
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
				"activeRental": bson.M{
					"$ifNull": []interface{}{
						bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{
									"$filter": bson.M{
										"input": "$rentals",
										"as":    "rental",
										"cond":  bson.M{"$eq": []interface{}{"$$rental.active", true}},
									},
								},
								0,
							},
						},
						bson.M{},
					},
				},
			},
		},
		{
			"$unwind": "$activeRental", // Entfalte das "activeRental"-Objekt
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "activeRental.userId",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$group": bson.M{
				"_id":          "$_id",
				"description":  bson.M{"$first": "$description"},
				"location":     bson.M{"$first": "$location"},
				"name":         bson.M{"$first": "$name"},
				"tagNames":     bson.M{"$push": "$tags.name"}, // Extrahiere die Tag-Namen in ein Array
				"activeRental": bson.M{"$first": "$activeRental"},
			},
		},
		{
			"$project": bson.M{
				"_id":         0, // Ausblenden der _id-Felder
				"description": 1,
				"location":    1,
				"name":        1,
				"tagNames":    1, // Das Array mit Tag-Namen beibehalten
				"activeRental": bson.M{
					"active": 1,
					"userId": "$activeRental.userId", // Die userId aus dem verknüpften "users"-Dokument verwenden
					"email":  "$user.email",          // Die E-Mail-Adresse aus dem verknüpften "users"-Dokument verwenden
				},
			},
		},
	}

	documents, err := database.NewDBAggregation(collection, pipeline)

	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	document := documents[0]

	c.IndentedJSON(200, document)
}
