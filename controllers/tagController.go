package controllers

import (
	"BorrowBox/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTagsById(itemId string) []string {
	tags := []string{}
	tagIds, err := database.GetDocumentsByCollectionFiltered("itemTag", "itemId", itemId, true) //TODO: ersetzen durch getDocumentsByCollectionFiltered
	if err != nil {
		return nil
	}
	for _, tagIdMapping := range tagIds {
		tagId := tagIdMapping["tagId"]
		if oid, ok := tagId.(primitive.ObjectID); ok {
			tagIdString := oid.Hex()
			tag, err := database.GetDocumentByID("tags", tagIdString)
			if err != nil {
				return nil
			}
			if tagName, ok := tag["name"].(string); ok {
				tags = append(tags, tagName)
			}
		}

	}
	return tags
}
