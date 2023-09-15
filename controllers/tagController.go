package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrCreateTags(tagNames []string) ([]primitive.ObjectID, error) {
	client, err := database.NewMongoDB()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("borrowbox").Collection("tags")
	var tagIDs []primitive.ObjectID

	// Erstellen Sie eine Liste der Tag-Namen, die gesucht werden sollen
	filter := bson.M{"name": bson.M{"$in": tagNames}}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Erstellen Sie eine Map, um vorhandene Tags nach Namen abzubilden
	existingTags := make(map[string]models.Tag)

	// Iterieren Sie über die gefundenen Tags und fügen Sie sie der Map hinzu
	for cursor.Next(context.Background()) {
		var tag models.Tag
		if err := cursor.Decode(&tag); err != nil {
			return nil, err
		}
		existingTags[tag.Name] = tag
	}

	// Erstellen Sie eine Liste von Tags zum Einfügen
	var tagsToInsert []interface{}

	// Iterieren Sie über die Tag-Namen aus dem Frontend
	for _, tagName := range tagNames {
		if tag, ok := existingTags[tagName]; ok {
			// Der Tag existiert bereits, fügen Sie ihn zur Liste hinzu
			tagIDs = append(tagIDs, tag.ID)
		} else {
			// Der Tag existiert nicht, fügen Sie ihn zur Liste hinzu
			tagsToInsert = append(tagsToInsert, bson.M{"_id": primitive.NewObjectID(), "name": tagName})
		}
	}

	if len(tagsToInsert) > 0 {
		results, err := collection.InsertMany(context.Background(), tagsToInsert)
		if err != nil {
			return nil, err
		}

		for _, result := range results.InsertedIDs {
			tagIDs = append(tagIDs, result.(primitive.ObjectID))
		}

	}

	return tagIDs, nil
}

func InsertTagItem(itemId primitive.ObjectID, tagIds []primitive.ObjectID) error {
	client, err := database.NewMongoDB()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("borrowbox").Collection("itemTag")
	// Erstellen Sie ein Array von Einträgen, wobei jede Tag-ID mit derselben Item-ID verknüpft ist
	var entries []interface{}
	for _, tagID := range tagIds {
		entry := bson.M{
			"itemId": itemId, // Die ID des Items
			"tagId":  tagID,  // Die ID des Tags
		}
		entries = append(entries, entry)
	}

	// Fügen Sie alle Einträge auf einmal in die Tabelle ein
	_, err = collection.InsertMany(context.Background(), entries)
	if err != nil {
		// Handle Fehler, falls die Einfügeoperation fehlschlägt
		return err
	}
	return nil
}
