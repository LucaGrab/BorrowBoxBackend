package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewMongoDB creates a new MongoDB instance.
func NewMongoDB() (*mongo.Client, error) {
	mongodbURI := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbURI))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	return client, nil
}

// InsertDocument fügt ein Dokument in die Sammlung ein.
func InsertDocument(collectionName string, document interface{}) error {
	client, err := NewMongoDB()
	if err != nil {
		return err
	}
	collection := client.Database("borrowbox").Collection(collectionName)
	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDocument aktualisiert ein Dokument in der Sammlung.
func UpdateDocument(collectionName string, documentID string, update interface{}) error {
	client, err := NewMongoDB()
	if err != nil {
		return err
	}
	collection := client.Database("borrowbox").Collection(collectionName)
	id, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDocument löscht ein Dokument aus der Sammlung.
func DeleteDocument(collectionName string, documentID string) error {
	client, err := NewMongoDB()
	if err != nil {
		return err
	}
	collection := client.Database("borrowbox").Collection(collectionName)
	id, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

// FindOne retrieves a single document from the collection based on a filter.
func getDocumentByID(collectionName string, documentID string) (bson.M, error) {
	client, err := NewMongoDB()
	collection := client.Database("borrowbox").Collection(collectionName)
	// ID in ObjectID konvertieren
	id, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return nil, err
	}
	// Query erstellen
	filter := bson.M{"_id": id}
	fmt.Println(documentID)
	// Ergebnis abrufen
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	fmt.Println("test")
	fmt.Println(result)

	return result, nil
}

func getAllDcoumentsByCollection(collectionName string) ([]bson.M, error) {
	client, err := NewMongoDB()
	collection := client.Database("borrowbox").Collection(collectionName)

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
	return results, nil
}
