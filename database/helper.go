package database

import (
	"BorrowBox/models"
	"bytes"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io"
	"mime/multipart"
	"os"
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
func InsertDocument(collectionName string, document interface{}) (primitive.ObjectID, error) {
	client, err := NewMongoDB()
	if err != nil {
		defer client.Disconnect(context.TODO())
		return primitive.NilObjectID, err
	}
	collection := client.Database("borrowbox").Collection(collectionName)
	result, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		defer client.Disconnect(context.TODO())
		return primitive.NilObjectID, err
	}
	defer client.Disconnect(context.TODO())
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		defer client.Disconnect(context.TODO())
		return primitive.NilObjectID, err
	}
	return insertedID, nil
}

// UpdateDocument aktualisiert ein Dokument in der Sammlung.
func UpdateDocument(collectionName string, documentID string, update interface{}) error {
	client, err := NewMongoDB()
	if err != nil {
		defer client.Disconnect(context.TODO())
		return err
	}
	collection := client.Database("borrowbox").Collection(collectionName)
	id, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		println("id error")
		defer client.Disconnect(context.TODO())
		return err
	}
	filter := bson.M{"_id": id}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		println("update error")
		defer client.Disconnect(context.TODO())
		return err
	}

	return nil
}

// DeleteDocument löscht ein Dokument aus der Sammlung.
func DeleteDocument(collectionName string, documentID string) error {
	client, err := NewMongoDB()
	if err != nil {
		defer client.Disconnect(context.TODO())
		return err
	}
	collection := client.Database("borrowbox").Collection(collectionName)
	id, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		defer client.Disconnect(context.TODO())
		return err
	}
	filter := bson.M{"_id": id}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		defer client.Disconnect(context.TODO())
		return err
	}
	return nil
}

// Lösche alle Dokumente mit einem bestimmten Filter
func DeleteAllDocuments(collectionName string, filter bson.M) error {
	client, err := NewMongoDB()
	if err != nil {
		defer client.Disconnect(context.Background())
		return err
	}
	collection := client.Database("borrowbox").Collection(collectionName)
	_, err = collection.DeleteMany(context.Background(), filter)
	if err != nil {
		defer client.Disconnect(context.Background())
		return err
	}
	return nil
}

// FindOne retrieves a single document from the collection based on a filter.
func GetDocumentByID(collectionName string, documentID string) (bson.M, error) {
	client, err := NewMongoDB()
	collection := client.Database("borrowbox").Collection(collectionName)
	// ID in ObjectID konvertieren
	id, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		defer client.Disconnect(context.TODO())
		return nil, err
	}
	// Query erstellen
	filter := bson.M{"_id": id}
	// Ergebnis abrufen
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		defer client.Disconnect(context.TODO())
		return nil, err
	}
	return result, nil
}

func GetAllDcoumentsByCollection(collectionName string) ([]bson.M, error) {
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
	defer client.Disconnect(context.TODO())
	return results, nil
}

func NewDBAggregation(collectionName string, pipeline []bson.M) ([]bson.M, error) {
	client, err := NewMongoDB()
	collection := client.Database("borrowbox").Collection(collectionName)
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	return results, nil
}

func GetActiveRentalByItemId(itemId primitive.ObjectID) models.RentalWithId {
	client, err := NewMongoDB()
	collection := client.Database("borrowbox").Collection("rentals")
	filter := bson.M{"itemId": itemId, "active": true}
	var result models.RentalWithId
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	return result
}

func SetItemImage(itemId string, image *multipart.FileHeader) {
	client, err := NewMongoDB()
	if err != nil {
		defer client.Disconnect(context.Background())
		return
	}

	db := client.Database("borrowbox")

	fs, err := gridfs.NewBucket(db)
	if err != nil {
		println("gridfs error")
		panic(err)
	}

	src, err := image.Open()
	if err != nil {
		println("open error")
		panic(err)
	}

	uploadStream, err := fs.OpenUploadStream(itemId)
	if err != nil {
		println("upload error")
		panic(err)
	}

	defer uploadStream.Close()

	_, err = io.Copy(uploadStream, src)
	if err != nil {
		println("copy error")
		panic(err)
	}
}

func GetItemImage(itemId string) []byte {
	client, err := NewMongoDB()
	if err != nil {
		defer client.Disconnect(context.Background())
		return nil
	}

	db := client.Database("borrowbox")

	fs, err := gridfs.NewBucket(db)
	if err != nil {
		println("gridfs error")
		panic(err)
	}

	var buf bytes.Buffer
	downloadStream, err := fs.OpenDownloadStreamByName(itemId)
	if err != nil {
		println("download error")
		return nil
	}
	defer downloadStream.Close()

	_, err = io.Copy(&buf, downloadStream)

	return buf.Bytes()
}
