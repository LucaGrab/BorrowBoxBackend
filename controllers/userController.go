package controllers

import (
	"BorrowBox/database"
	"BorrowBox/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	client, err := database.NewMongoDB()
	collection := client.Database("borrowbox").Collection("users")

	// Query erstellen
	filter := bson.D{} // Hier kannst du optional eine Filterbedingung angeben

	// Ergebnisse abrufen
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}

		// Hier kannst du die Projektion auf die gewünschten Felder anwenden
		projectionResult := bson.M{
			"id":       result["_id"],
			"role":     result["role"],
			"username": result["username"],
			"email":    result["email"],
		}

		results = append(results, projectionResult)
	}

	if err := cursor.Err(); err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	c.IndentedJSON(200, results)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := database.DeleteDocument("users", id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "User deleted"})
}

func UserById(c *gin.Context) {
	id := c.Param("id")
	user, err := database.GetDocumentByID("users", id)
	if err != nil {
		c.IndentedJSON(404, gin.H{"message": err.Error()})
		return
	}

	// Entfernen Sie das Passwort aus dem Benutzerobjekt
	delete(user, "password")
	fmt.Println(user)
	c.IndentedJSON(200, user)
}

func InsertUser(c *gin.Context) {
	var newUser models.User // Ersetze YourDataStruct mit der tatsächlichen Struktur deiner Daten

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	_, err := database.InsertDocument("users", newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "User inserted successfully"})
}

func UpdateUser(c *gin.Context) {
	var updatedUser models.User

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hier erstellst du eine Filterbedingung, um den Benutzer in der MongoDB zu identifizieren (z. B. anhand der ID).
	filter := bson.M{"_id": updatedUser.ID}

	// Hier erstellst du eine Aktualisierungsanweisung, um die Benutzerdaten zu aktualisieren.
	update := bson.M{
		"$set": bson.M{},
	}

	// Überprüfe, ob das Email-Feld im JSON-Parameter gefüllt ist, bevor du es aktualisierst.
	if updatedUser.Email != "" {
		update["$set"].(bson.M)["email"] = updatedUser.Email
	}
	fmt.Println(updatedUser.Password)
	// Überprüfe, ob das Password-Feld im JSON-Parameter gefüllt ist, bevor du es aktualisierst.
	if updatedUser.Password != "" {
		update["$set"].(bson.M)["password"] = updatedUser.Password
	}

	// Überprüfe, ob das Username-Feld im JSON-Parameter gefüllt ist, bevor du es aktualisierst.
	if updatedUser.Username != "" {
		update["$set"].(bson.M)["username"] = updatedUser.Username
	}

	client, err := database.NewMongoDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Öffne eine Verbindung zur MongoDB und aktualisiere den Benutzer.
	collection := client.Database("borrowbox").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Benutzeraktualisierung fehlgeschlagen"})
		return
	}

	// Erfolgreiche Aktualisierung
	c.JSON(http.StatusOK, gin.H{"message": "Benutzer erfolgreich aktualisiert"})
}
