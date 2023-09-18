package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tag struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `json:"name"`
}

type UserTag struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	userId primitive.ObjectID `json:"userId"`
	tagId  primitive.ObjectID `json:"tagId"`
}
