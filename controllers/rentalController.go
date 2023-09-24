package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"context"
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

	//wenn die report beschreibung nicht leer ist, dann report in die tabelle einfügen
	if returnData.ReportDescription != "" {
		println("reporting item...")
		report := models.Report{
			ItemId:        itemId,
			Time:          rental.End,
			UserId:        returnData.UserId,
			Description:   returnData.ReportDescription,
			StateCritical: returnData.ReportStateCritical,
		}
		_, err := database.InsertDocument("reports", report)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert report"})
			return
		}
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Rental updated successfully"})
}

func GetHistory(c *gin.Context) {
	// User ID aus dem Pfadparameter abrufen
	userID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		// Handle den Fehler hier, falls die Konvertierung fehlschlägt
	}
	// Benutzerrolle abrufen (z. B. aus Ihrer Datenbank)
	userRole, err := getUserRole(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Benutzerrolle"})
		return
	}

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "items",  // Name der zweiten Sammlung
				"localField":   "itemId", // Feld in der ersten Sammlung
				"foreignField": "_id",    // Feld in der zweiten Sammlung
				"as":           "item",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "users",  // Name der dritten Sammlung
				"localField":   "userId", // Feld in der ersten Sammlung
				"foreignField": "_id",    // Feld in der dritten Sammlung
				"as":           "user",
			},
		},
		{
			"$unwind": "$item", // Verflachen des "item"-Arrays
		},
		{
			"$unwind": "$user", // Verflachen des "user"-Arrays
		},
		{
			"$match": bson.M{
				"user._id": id, // Filter nach der gewünschten UserID
			},
		},
		{
			"$project": bson.M{
				"_id":      0, // Das "_id"-Feld ausblenden
				"start":    1,
				"end":      1,
				"active":   1,
				"userName": "$user.username", // "username" aus dem verknüpften "user"-Dokument
				"userId":   "$user._id",
				"itemName": "$item.name", // "name" aus dem verknüpften "item"-Dokument
			},
		},
	}

	// Wenn der Benutzer ein Admin ist, alle Rentals abrufen
	if userRole == "admin" {
		pipeline = []bson.M{
			{
				"$lookup": bson.M{
					"from":         "items",  // Name der zweiten Sammlung
					"localField":   "itemId", // Feld in der ersten Sammlung
					"foreignField": "_id",    // Feld in der zweiten Sammlung
					"as":           "item",
				},
			},
			{
				"$lookup": bson.M{
					"from":         "users",  // Name der dritten Sammlung
					"localField":   "userId", // Feld in der ersten Sammlung
					"foreignField": "_id",    // Feld in der dritten Sammlung
					"as":           "user",
				},
			},
			{
				"$unwind": "$item", // Verflachen des "item"-Arrays
			},
			{
				"$unwind": "$user", // Verflachen des "user"-Arrays
			},
			{
				"$project": bson.M{
					"_id":      0, // Das "_id"-Feld ausblenden
					"start":    1,
					"end":      1,
					"active":   1,
					"userName": "$user.username", // "username" aus dem verknüpften "user"-Dokument
					"userId":   "$user._id",
					"itemName": "$item.name", // "name" aus dem verknüpften "item"-Dokument
				},
			},
		}
	}

	sortStage := bson.M{
		"$sort": bson.M{
			"start": -1,
		},
	}

	// Pipeline um die Sortierstufe erweitern
	pipeline = append(pipeline, sortStage)

	// Aggregation durchführen
	rentals, err := database.NewDBAggregation("rentals", pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Daten"})
		return
	}

	// JSON-Array als Antwort senden
	c.JSON(http.StatusOK, rentals)
}

// getUserRole ruft die Benutzerrolle (Admin oder nicht) basierend auf der Benutzer-ID aus der Datenbank ab.
func getUserRole(userID primitive.ObjectID) (string, error) {
	// Verbinden Sie sich mit Ihrer Datenbank
	client, err := database.NewMongoDB()
	if err != nil {
		return "", err
	}
	defer client.Disconnect(context.Background())

	// Rufen Sie die Benutzerrolle aus der Datenbank ab
	collection := client.Database("borrowbox").Collection("users")
	filter := bson.M{"_id": userID}
	var user struct {
		Role string `bson:"role"`
	}
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return "", err
	}

	return user.Role, nil
}
