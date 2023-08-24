package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://root:revxe3-roxfUb-wepcih@cluster1.epptnkq.mongodb.net/"))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Verbindung zur MongoDB hergestellt!")
	collection := client.Database("borrowbox").Collection("users")
	fmt.Println(collection)

	// Query erstellen
	filter := bson.D{} // Hier kannst du optional eine Filterbedingung angeben

	// Ergebnisse abrufen
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())

	// Ergebnisse verarbeiten
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		panic(err)
	}

	// Ergebnisse ausgeben
	for _, result := range results {

		fmt.Println(result)
	}

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"test":    results,
			"message": "Hello, World!",
		})
	})
	r.Run(":8080") // Starte den Gin-Server auf Port 8080
}
