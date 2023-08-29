package main

import (
	"context"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Id    string `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

func userById(c *gin.Context) {
	id := c.Param("id")
	user, err := getUserById(id)
	if err != nil {
		return
	}
	c.IndentedJSON(200, user)
}

func getUserById(id string) (User, error) {
	/*userCollection, err := NewMongoDB("mongodb+srv://root:revxe3-roxfUb-wepcih@cluster1.epptnkq.mongodb.net/", "borrowbox", "users")
	filter := bson.D{{"id", id}}
	User := userCollection.FindOne(filter)*/

	return User{}, nil
}

func getUsers(c *gin.Context) {

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

	c.IndentedJSON(200, results)

	/*userCollection, err := NewMongoDB("mongodb+srv://root:revxe3-roxfUb-wepcih@cluster1.epptnkq.mongodb.net/", "borrowbox", "users")
	users := userCollection.GetAll()
	if err != nil {
		return
	}
	c.IndentedJSON(200, users)*/
}

func startGinServer() {

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8100"} // Add your frontend addresses here
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	r.Use(cors.New(config))

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"test":    "results",
			"message": "Hello, World!",
		})
	})

	r.GET("users/:id", userById)
	r.GET("users", getUsers)

	r.Run(":8080") // Starte den Gin-Server auf Port 8080
}
