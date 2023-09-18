package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ItemTagForInsert struct {
	ItemId primitive.ObjectID `json:"itemId" bson:"itemId"`
	TagId  primitive.ObjectID `json:"tagId" bson:"tagId"`
}
