package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ...
func GetTags(c *gin.Context) {
	tags, err := database.GetAllDcoumentsByCollection("tags")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tags"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

func GetAllTags(c *gin.Context) {
	// Holen Sie sich die userId aus dem Parameter
	userID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	// Verbinden Sie sich mit der MongoDB und wählen Sie die "tags"-Sammlung aus
	client, err := database.NewMongoDB()
	if err != nil {
		c.JSON(500, gin.H{"message": "Database error"})
		return
	}
	tagsCollection := client.Database("borrowbox").Collection("tags")

	// Durchführen einer Find-Abfrage in der "tags"-Sammlung, um alle Tags abzurufen
	tagsCursor, err := tagsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"message": "Database error"})
		return
	}
	defer tagsCursor.Close(context.TODO())

	// Extrahieren Sie die gefundenen Tags aus dem Tags-Cursor
	var tags []bson.M
	if err := tagsCursor.All(context.TODO(), &tags); err != nil {
		c.JSON(500, gin.H{"message": "Error while decoding tags"})
		return
	}

	// Verbinden Sie sich mit der MongoDB und wählen Sie die "userTags"-Sammlung aus
	userTagsCollection := client.Database("borrowbox").Collection("userTags")

	// Erstellen Sie eine Filterbedingung für die userId
	userTagFilter := bson.M{"userId": userID}

	// Durchführen einer Find-Abfrage in der "userTags"-Sammlung mit dem Filter
	userTagsCursor, err := userTagsCollection.Find(context.TODO(), userTagFilter)
	if err != nil {
		c.JSON(500, gin.H{"message": "Database error"})
		return
	}
	defer userTagsCursor.Close(context.TODO())

	// Extrahieren Sie die gefundenen userTags aus dem Cursor
	var userTagDocs []bson.M
	if err := userTagsCursor.All(context.TODO(), &userTagDocs); err != nil {
		c.JSON(500, gin.H{"message": "Error while decoding userTags"})
		return
	}

	// Extrahieren Sie die tagIds aus den userTag-Dokumenten
	var userTagIDs []primitive.ObjectID
	for _, doc := range userTagDocs {
		tagID, ok := doc["tagId"].(primitive.ObjectID)
		if ok {
			userTagIDs = append(userTagIDs, tagID)
		}
	}

	// Markieren Sie die Tags, die zu einem Benutzer gehören
	for i, tag := range tags {
		tagID, ok := tag["_id"].(primitive.ObjectID)
		if ok && contains(userTagIDs, tagID) {
			tags[i]["tagged"] = true
		} else {
			tags[i]["tagged"] = false
		}
	}

	c.JSON(200, tags)
}

// Hilfsfunktion zur Überprüfung, ob ein Wert in einem Slice vorhanden ist
func contains(slice []primitive.ObjectID, item primitive.ObjectID) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func UpdateUserTag(c *gin.Context) {
	var userTagData map[string]string
	if err := c.ShouldBindJSON(&userTagData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	client, err := database.NewMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	tags := client.Database("borrowbox").Collection("userTags")

	userID, err := primitive.ObjectIDFromHex(userTagData["userId"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	tagID, err := primitive.ObjectIDFromHex(userTagData["tagId"])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tagId"})
		return
	}

	userTagDoc := bson.M{
		"userId": userID,
		"tagId":  tagID,
	}
	var existingUserTag models.UserTag
	err = tags.FindOne(context.TODO(), userTagDoc).Decode(&existingUserTag)
	if err == nil {
		_, err := tags.DeleteOne(context.TODO(), userTagDoc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting document"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "UserTag deleted"})
		return
	} else if err == mongo.ErrNoDocuments {
		fmt.Println(userTagDoc)
		_, err := tags.InsertOne(context.TODO(), userTagDoc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting document"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "UserTag added"})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while querying the database"})
		return
	}
}

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

func CreateTag(c *gin.Context) {
	// Struktur zur Analyse des JSON-Objekts definieren
	var tagData struct {
		LoginToken string `json:"loginToken"`
		FilterName string `json:"filterName"`
	}

	// JSON-Daten analysieren
	if err := c.ShouldBindJSON(&tagData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Verbindung zur MongoDB herstellen
	client, err := database.NewMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer client.Disconnect(context.TODO())

	// MongoDB-Sammlung "tags" auswählen
	tagsCollection := client.Database("borrowbox").Collection("tags")

	// Überprüfen, ob der FilterName bereits in der Sammlung vorhanden ist
	filter := bson.M{"name": tagData.FilterName}
	var existingTag models.Tag
	err = tagsCollection.FindOne(context.TODO(), filter).Decode(&existingTag)
	if err == nil {
		// Der Tag existiert bereits, eine entsprechende Fehlermeldung zurückgeben
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag with the same name already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		// Ein anderer Fehler ist aufgetreten, eine Fehlermeldung zurückgeben
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while querying the database"})
		return
	}

	// Den Tag in die Sammlung einfügen, da er nicht vorhanden ist
	tagID := primitive.NewObjectID()
	newTag := models.Tag{
		ID:   tagID,
		Name: tagData.FilterName,
	}

	_, err = tagsCollection.InsertOne(context.TODO(), newTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting the new tag"})
		return
	}

	// Aktualisierte Liste aller Tags abrufen und zurücksenden
	tags, err := database.GetAllDcoumentsByCollection("tags")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tags"})
		return
	}

	// Erfolgreiche Antwort mit der aktualisierten Liste aller Tags zurückgeben
	c.JSON(http.StatusOK, tags)
}

func DeleteTag(c *gin.Context) {
	// Hier nehmen wir an, dass die ID im Request Body als JSON-Objekt übergeben wird.
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Überprüfen Sie, ob die ID im Request Body vorhanden ist.
	id, ok := requestBody["id"].(string)
	if !ok || id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid or missing 'id' in request body"})
		return
	}

	// Führen Sie die Löschoperation mit der extrahierten ID durch.
	err := database.DeleteDocument("tags", id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// Aktualisierte Liste aller Tags abrufen und zurücksenden
	tags, err := database.GetAllDcoumentsByCollection("tags")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tags"})
		return
	}

	// Erfolgreiche Antwort mit der aktualisierten Liste aller Tags zurückgeben
	c.JSON(http.StatusOK, tags)
}
