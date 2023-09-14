package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	TagNames      []string           `bson:"tagNames" json:"tagNames"`
	Name          string             `json:"name"`
	Location      string             `json:"location"`
	Description   string             `json:"description"`
	Available     bool               `json:"available"`
	CurrentRenter string             `json:"currentRenter"`
}
