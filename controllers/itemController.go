package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteItem(c *gin.Context) {
	id := c.Param("id")

	// Erstelle ein Update-Filter, um das "deleted" Feld auf true zu setzen
	update := bson.M{"$set": bson.M{"deleted": true}}

	// Führe das Update in der "items"-Tabelle aus
	err := database.UpdateDocument("items", id, update)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Item marked as deleted successfully"})
}

func UpdateItem(c *gin.Context) {
	var updatedItem models.ItemMitTagIds
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	updateDoc := bson.M{
		"$set": bson.M{
			"name":        updatedItem.Name,
			"location":    updatedItem.Location,
			"description": updatedItem.Description,
		},
	}

	err := database.UpdateDocument("items", updatedItem.ID.Hex(), updateDoc)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}
	//das brauche ich wahrscheinlich nichtmehr weil ids im frontend bekannt sind -
	//vllt aber doch um zu schauen dass tags nicht in der zwischenzeit hinzugefügt wurden
	/*
		tags, err := GetOrCreateTags(updatedItem.TagNames)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error - Failed to insert tags"})
			return
		}*/
	err = UpdateItemTags(updatedItem.ID, updatedItem.TagIds)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error - Failed to insert item tag mapping"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Item updated successfully"})

}

func UploadItemImage(c *gin.Context) {
	itemId := c.Param("id")

	file, err := c.FormFile("photo")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	database.SetItemImage(itemId, file)
	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded and saved successfully"})
}

func InsertItem(c *gin.Context) {

	var newItem models.AddItem

	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	if newItem.Name == "NEIN" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	newItem.ID = primitive.NewObjectID()
	//item ohne tags haben
	itemForInsert := models.ItemForInsert{
		ID:          newItem.ID,
		Name:        newItem.Name,
		Location:    newItem.Location,
		Description: newItem.Description,
	}
	_, err := database.InsertDocument("items", itemForInsert)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert item"})
		return
	}
	//das brauche ich wahrscheinlich nichtmehr weil ids im frontend bekannt sind -
	//vllt aber doch um zu schauen dass tags nicht in der zwischenzeit hinzugefügt wurden
	/*
		tags, err := GetOrCreateTags(newItem.TagNames)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error - Failed to insert tags"})
			return
		}*/
	err = InsertTagItem(newItem.ID, newItem.TagIds)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error - Failed to insert item tag mapping"})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Item inserted successfully"})
}

func GetItems(c *gin.Context) {
	collection := "items"
	pipeline := []bson.M{
		{
			"$match": bson.M{"deleted": false}, // Filtere nach "deleted: false"
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
			"$lookup": bson.M{
				"from": "reports",
				"let":  bson.M{"itemId": "$_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{"$eq": []interface{}{"$itemId", "$$itemId"}},
						},
					},
					{
						"$sort": bson.M{"time": -1},
					},
					{
						"$limit": 1,
					},
				},
				"as": "latestReport",
			},
		},
		{
			"$addFields": bson.M{
				"latestReport": bson.M{"$arrayElemAt": []interface{}{"$latestReport", 0}},
			},
		},
		{
			"$addFields": bson.M{
				"available": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": []interface{}{"$latestReport.statecritical", true}},
						"then": false,
						"else": "$available",
					},
				},
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
		c.IndentedJSON(400, gin.H{"message": "Invalid ID"})
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
			"$unwind": "$items",
		},
		{
			"$match": bson.M{"items.deleted": false}, // Filtere nach "deleted: false"
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

	if len(documents) == 0 {
		c.IndentedJSON(http.StatusNoContent, gin.H{"message": "No matching documents found"})
		return
	}

	document := documents[0]
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
					"active":    1,
					"userId":    "$activeRental.userId", // Die userId aus dem verknüpften "users"-Dokument verwenden
					"email":     "$user.email",
					"startTime": "$activeRental.start", // Die E-Mail-Adresse aus dem verknüpften "users"-Dokument verwenden
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

	var item models.ItemMitReport
	item = models.ItemMitReport{
		ID:          document["_id"].(primitive.ObjectID),
		TagNames:    tagNamesSlice,
		Name:        document["name"].(string),
		Location:    document["location"].(string),
		Description: document["description"].(string),
	}
	if document["activeRental"].(primitive.M)["startTime"] != nil {
		startTimePrimitive := document["activeRental"].(primitive.M)["startTime"].(primitive.DateTime)
		startTimeTime := startTimePrimitive.Time()
		startTimeFormatted := startTimeTime.Format("2006-01-02 15:04")
		item.RentedSince = startTimeFormatted
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

	//---------------- den report seperat holen
	pipeline2 := []bson.M{
		{
			"$match": bson.M{
				"itemId": formattedId, // Filtern nach der Item-ID
			},
		},
		{
			"$sort": bson.M{
				"time": -1, // Sortieren nach "time" absteigend, um den neuesten Bericht zuerst zu erhalten
			},
		},
		{
			"$limit": 1, // Begrenzen auf den neuesten Bericht
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "userId",      // Das Feld in der "reports" Tabelle
				"foreignField": "_id",         // Das Feld in der "users" Tabelle
				"as":           "userDetails", // Das Alias für das Ergebnis
			},
		},
		{
			"$addFields": bson.M{
				"username": bson.M{"$arrayElemAt": []interface{}{"$userDetails.username", 0}},
			},
		},
		{
			"$project": bson.M{
				"_id":           0, // Ausblenden des _id-Feldes
				"itemId":        1,
				"time":          1,
				"description":   1,
				"username":      1,
				"statecritical": 1,

				// Fügen Sie hier weitere Felder aus dem "reports"-Dokument hinzu, wenn benötigt
			},
		},
	}

	reports, err := database.NewDBAggregation("reports", pipeline2)

	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}
	if len(reports) > 0 {
		report := reports[0]
		item.ReportDescription = report["description"].(string)
		item.ReportTime = report["time"].(primitive.DateTime).Time().Format("2006-01-02 15:04")
		item.ReportStateCritical = report["statecritical"].(bool)
		item.ReportUser = report["username"].(string)
		fmt.Println(report["statecritical"].(bool))
		if report["statecritical"].(bool) {
			item.Available = false
		}
	}

	//----------------

	//-------------------------------------------------------------------------------------------
	c.IndentedJSON(200, item)
}
