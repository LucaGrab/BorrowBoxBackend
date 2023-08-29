package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB wraps the MongoDB client and establishes connections.
type MongoDB struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// NewMongoDB creates a new MongoDB instance.
func NewMongoDB(connectionString, collectionName, databaseName string) (*MongoDB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	return &MongoDB{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

// Close disconnects from the MongoDB server.
func (db *MongoDB) Close() {
	err := db.client.Disconnect(context.Background())
	if err != nil {
		log.Println("Error disconnecting from MongoDB:", err)
	}
}

// InsertOne inserts a document into the collection.
func (db *MongoDB) InsertOne(document interface{}) error {
	_, err := db.collection.InsertOne(context.Background(), document)
	return err
}

// FindOne retrieves a single document from the collection based on a filter.
func (db *MongoDB) FindOne(filter interface{}) error {
	var result []bson.M
	return db.collection.FindOne(context.Background(), filter).Decode(result)
}

// UpdateOne updates a document in the collection based on a filter.
func (db *MongoDB) UpdateOne(filter interface{}, update interface{}) error {
	_, err := db.collection.UpdateOne(context.Background(), filter, update)
	return err
}

// DeleteOne deletes a document from the collection based on a filter.
func (db *MongoDB) DeleteOne(filter interface{}) error {
	_, err := db.collection.DeleteOne(context.Background(), filter)
	return err
}

// GetAll retrieves all documents from a collection.
func (db *MongoDB) GetAll() error {
	// Query erstellen
	filter := bson.D{} // Hier kannst du optional eine Filterbedingung angeben

	// Ergebnisse abrufen
	cursor, err := db.collection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())

	// Ergebnisse verarbeiten
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		panic(err)
	}
	fmt.Println("results")
	// Ergebnisse ausgeben
	for _, result := range results {

		fmt.Println(result)
	}
	return nil

	/*fmt.Println("test")

	var results []bson.M
	cur, err := db.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	fmt.Println("test")

	err = cur.All(context.Background(), results)
	if err != nil {
		return err
	}
	fmt.Println(results)
	fmt.Println("test")
	for _, result := range results {

		fmt.Println(result)
	}
	return nil*/

}
