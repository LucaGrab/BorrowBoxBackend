package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	Time "time"

	"github.com/gin-gonic/gin"
)

func InsertReport(c *gin.Context) {
	var report models.Report
	report.Time = Time.Now()
	c.BindJSON(&report)
	fmt.Println(report)
	_, err := database.InsertDocument("reports", report)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Report could not be inserted.",
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Report inserted.",
		})
	}
}

func GetReports(c *gin.Context) {
	pipeline := []bson.M{
		{"$match": bson.M{"statecritical": true}},

		// Fügen Sie ein Lookup für das "items"-Dokument hinzu
		{
			"$lookup": bson.M{
				"from":         "items", // Name der "items"-Collection
				"localField":   "itemId",
				"foreignField": "_id",
				"as":           "itemData",
			},
		},

		// Fügen Sie ein Lookup für das "users"-Dokument hinzu
		{
			"$lookup": bson.M{
				"from":         "users", // Name der "users"-Collection
				"localField":   "userId",
				"foreignField": "_id",
				"as":           "userData",
			},
		},

		// Projizieren Sie nur die gewünschten Felder
		{
			"$project": bson.M{
				"description": 1,
				"time":        1,
				"itemName":    bson.M{"$arrayElemAt": []interface{}{"$itemData.name", 0}},
				"userName":    bson.M{"$arrayElemAt": []interface{}{"$userData.username", 0}},
				"itemId":      bson.M{"$arrayElemAt": []interface{}{"$itemData._id", 0}},
			},
		},
		{
			"$unwind": "$itemName", // "itemName" entpacken
		},
		{
			"$unwind": "$userName", // "userName" entpacken
		},
	}

	sortStage := bson.M{
		"$sort": bson.M{
			"time": -1,
		},
	}

	// Pipeline um die Sortierstufe erweitern
	pipeline = append(pipeline, sortStage)

	reports, err := database.NewDBAggregation("reports", pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Daten"})
		return
	}

	// JSON-Array als Antwort senden
	c.JSON(http.StatusOK, reports)
}
