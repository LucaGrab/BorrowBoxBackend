package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetItems(c *gin.Context) {
	collection := "items"
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "rentals",
				"localField":   "_id",
				"foreignField": "itemId",
				"as":           "rentals",
			},
		},
		{
			"$addFields": bson.M{
				"available": bson.M{
					"$not": bson.M{
						"$anyElementTrue": bson.M{
							"$map": bson.M{
								"input": "$rentals",
								"as":    "rental",
								"in":    "$$rental.active",
							},
						},
					},
				},
			},
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
				"from":         "tags",
				"localField":   "itemTags.tagId",
				"foreignField": "_id",
				"as":           "tags",
			},
		},
		{
			"$project": bson.M{
				"_id":         1,
				"description": 1,
				"location":    1,
				"name":        1,
				"tags": bson.M{
					"$map": bson.M{
						"input": "$tags",
						"as":    "tag",
						"in":    "$$tag.name",
					},
				},
				"available": 1,
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
			"$unwind": "$items", // Entfalte das "items"-Array
		},
		{
			"$group": bson.M{
				"_id":   nil,
				"items": bson.M{"$push": "$items"},
			},
		},
		{
			"$project": bson.M{
				"_id":   0,
				"items": 1,
			},
		},
	}
	documents, err := database.NewDBAggregation("rentals", pipeline)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	document := documents[0] // entfernt random array was sonst immer kommt - zeigt trotzdem mehr items
	c.IndentedJSON(200, document)
}

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
				"_id":         1, // Ausblenden der _id-Felder
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

	//------------------------- hier wird das document in ein item umgewandelt --------------------
	//tags müssen iwie umgewandelt werden
	var tagNamesSlice []string
	tagNamesPrimitiveA, ok := document["tagNames"].(primitive.A)
	if ok {
		for _, tag := range tagNamesPrimitiveA {
			tagNamesSlice = append(tagNamesSlice, tag.(string))
		}
	}

	var item models.Item
	item = models.Item{
		ID:          document["_id"].(primitive.ObjectID),
		TagNames:    tagNamesSlice,
		Name:        document["name"].(string),
		Location:    document["location"].(string),
		Description: document["description"].(string),
	}
	//weil das mit user join oben nicht geht
	activeRental, ok := document["activeRental"].(primitive.M)
	if ok {
		active, activeOK := activeRental["active"].(bool)
		if activeOK {
			item.Available = !active
			if active {
				fmt.Println(activeRental["userId"].(primitive.ObjectID).Hex())
				user, err := database.GetDocumentByID("users", activeRental["userId"].(primitive.ObjectID).Hex())
				if err != nil {
					c.IndentedJSON(404, gin.H{"message": err.Error()})
					return
				}
				item.CurrentRenter = user["username"].(string)
			}
		} else {
			item.Available = true
		}
	} else {
		item.Available = true
	}
	//-------------------------------------------------------------------------------------------
	c.IndentedJSON(200, item)
}
