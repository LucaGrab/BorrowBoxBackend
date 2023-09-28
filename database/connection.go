package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect() {
	mongodbURI := os.Getenv("MONGODB_URI")
	// Verbindung zur MongoDB herstellen
	clientOptions := options.Client().ApplyURI(mongodbURI) // Ersetzen Sie die URI durch die Ihrer MongoDB-Instanz.
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic("Konnte keine Verbindung zur MongoDB-Datenbank herstellen")
	}

	// Überprüfen Sie die Verbindung zur Datenbank
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic("Konnte die Verbindung zur MongoDB-Datenbank nicht überprüfen")
	}

}

// Schließen Sie die Verbindung zur Datenbank, wenn Ihr Programm beendet ist.
func CloseConnection() {
	err := DB.Client().Disconnect(context.TODO())
	if err != nil {
		panic("Fehler beim Schließen der Verbindung zur MongoDB-Datenbank")
	}
}
