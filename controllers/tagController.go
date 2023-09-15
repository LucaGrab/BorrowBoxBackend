package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrCreateTags(tagNames []string) ([]models.Tag, error) {
	client, err := database.NewMongoDB()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("borrowbox").Collection("tags")
	var tags []models.Tag

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
			tags = append(tags, tag)
		} else {
			// Der Tag existiert nicht, fügen Sie ihn zur Liste hinzu
			tagsToInsert = append(tagsToInsert, bson.M{"_id": primitive.NewObjectID(), "name": tagName})
		}
	}

	if len(tagsToInsert) > 0 {
		_, err := collection.InsertMany(context.Background(), tagsToInsert)
		if err != nil {
			return nil, err
		}

		// Fügen Sie alle Tags aus tagsToInsert zu tags hinzu
		for _, tagToInsert := range tagsToInsert {
			tag, ok := tagToInsert.(bson.M)
			if !ok {
				return nil, err
			}

			objectID, ok := tag["_id"].(primitive.ObjectID)
			if !ok {
				return nil, err
			}

			name, _ := tag["name"].(string)

			newTag := models.Tag{
				ID:   objectID,
				Name: name,
			}
			tags = append(tags, newTag)
		}
	}

	return tags, nil
}
